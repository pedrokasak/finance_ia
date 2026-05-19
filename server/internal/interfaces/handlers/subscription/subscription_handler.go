package subscription

import (
	"finance-ia/internal/domain/payment"
	"finance-ia/internal/domain/subscription"
	"finance-ia/internal/domain/user"
	dbsub "finance-ia/internal/infrastructure/database/subscription"
	paymentprovider "finance-ia/internal/infrastructure/payment/provider"
	"finance-ia/internal/utils"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionHandler struct {
	gateway          payment.PaymentGateway
	paymentProvider  string
	subscriptionRepo subscription.SubscriptionRepository
	planRepo         *dbsub.PlanRepository
	userRepo         user.Repository
	db               *gorm.DB
}

func NewSubscriptionHandler(
	gateway payment.PaymentGateway,
	paymentProvider string,
	subscriptionRepo subscription.SubscriptionRepository,
	planRepo *dbsub.PlanRepository,
	userRepo user.Repository,
	db *gorm.DB,
) *SubscriptionHandler {
	return &SubscriptionHandler{
		gateway:          gateway,
		paymentProvider:  paymentprovider.Normalize(paymentProvider),
		subscriptionRepo: subscriptionRepo,
		planRepo:         planRepo,
		userRepo:         userRepo,
		db:               db,
	}
}

func (h *SubscriptionHandler) RegisterRoutes(public, protected gin.IRouter) {
	sub := protected.Group("/subscription")
	{
		sub.GET("/", h.GetMySubscription)
		sub.GET("/plans", h.GetPlans)
		sub.GET("/invoices", h.GetInvoices)
		sub.POST("/checkout", h.CreateCheckout)
		sub.POST("/portal", h.CreatePortal)
	}

	public.POST("/webhook/stripe", h.HandleStripeWebhook)
	public.POST("/webhook/abacatepay", h.HandleAbacatePayWebhook)
}

// GetPlans returns all active plans from the database (with Stripe Price IDs)
func (h *SubscriptionHandler) GetPlans(c *gin.Context) {
	plans, err := h.planRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch plans"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"plans": plans})
}

func (h *SubscriptionHandler) GetMySubscription(c *gin.Context) {
	userID := utils.GetUserID(c)

	u, _ := h.userRepo.FindByID(userID)

	// Try local subscription record first
	sub, err := h.subscriptionRepo.FindByUserID(userID)
	if err == nil && sub != nil && sub.Plan != "" && sub.Plan != "free" {
		// We have a valid paid subscription record
		c.JSON(http.StatusOK, gin.H{
			"subscription": sub,
			"features":     subscription.GetPlanFeatures(sub.Plan),
		})
		return
	}

	// No local record (or free) — check if user has a Stripe customer with active subscription
	// This auto-heals when webhooks were missed during development
	if h.paymentProvider == paymentprovider.Stripe && u != nil && u.StripeCustomerID != "" {
		if stripeSub, serr := h.syncSubscriptionFromStripe(u); serr == nil && stripeSub != nil {
			c.JSON(http.StatusOK, gin.H{
				"subscription": stripeSub,
				"features":     subscription.GetPlanFeatures(stripeSub.Plan),
			})
			return
		}
	}

	// Fallback: read plan from users table (updated by webhooks)
	plan := "free"
	if u != nil && string(u.Plan) != "" {
		plan = string(u.Plan)
	}
	c.JSON(http.StatusOK, gin.H{
		"plan":     plan,
		"status":   "active",
		"features": subscription.GetPlanFeatures(plan),
	})
}

// syncSubscriptionFromStripe fetches active Stripe subscriptions for the user directly from Stripe.
// This auto-heals the local DB when webhooks were missed (e.g. wrong URL during development).
func (h *SubscriptionHandler) syncSubscriptionFromStripe(u *user.User) (*subscription.Subscription, error) {
	stripeSubs, err := h.gateway.ListActiveSubscriptions(u.StripeCustomerID)
	if err != nil {
		return nil, fmt.Errorf("stripe: list active subscriptions: %w", err)
	}
	if len(stripeSubs) == 0 {
		return nil, fmt.Errorf("no active subscription in stripe for customer %s", u.StripeCustomerID)
	}

	// Use the first active subscription
	stripeInfo := stripeSubs[0]

	// Resolve plan from PriceID
	planRecord, err := h.planRepo.FindByStripePriceID(stripeInfo.PriceID)
	if err != nil || planRecord == nil {
		// Could not resolve plan from price ID
		return nil, fmt.Errorf("stripe: could not resolve plan slug from price ID %s", stripeInfo.PriceID)
	}
	plan := planRecord.Slug

	// Upsert local subscription record
	sub, _ := h.subscriptionRepo.FindByUserID(u.ID)
	if sub == nil {
		sub = &subscription.Subscription{}
		sub.UserID = u.ID
	}
	sub.ExternalID = stripeInfo.ExternalID
	sub.Plan = plan
	sub.Status = subscription.StatusActive
	if err := h.subscriptionRepo.Upsert(sub); err != nil {
		return nil, fmt.Errorf("upsert subscription: %w", err)
	}

	// Persist plan on the user record so future fallbacks are immediate
	h.db.Model(&user.User{}).Where("id = ?", u.ID).Update("plan", plan)

	return sub, nil
}

func (h *SubscriptionHandler) GetInvoices(c *gin.Context) {
	userID := utils.GetUserID(c)

	u, err := h.userRepo.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if u.StripeCustomerID == "" {
		c.JSON(http.StatusOK, gin.H{"invoices": []interface{}{}})
		return
	}

	invoices, err := h.gateway.ListInvoices(u.StripeCustomerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list invoices"})
		return
	}

	if invoices == nil {
		invoices = []*payment.InvoiceInfo{}
	}

	c.JSON(http.StatusOK, gin.H{"invoices": invoices})
}

type createCheckoutRequest struct {
	Plan          string `json:"plan" binding:"required,oneof=premium pro"`
	BillingType   string `json:"billing_type" binding:"required,oneof=monthly yearly"`
	PaymentMethod string `json:"payment_method"` // "card" (default) | "pix" | "boleto"
}

func resolvePaymentMethods(method string) []string {
	switch method {
	case "pix":
		// PIX works on one-time payments; for subscriptions Stripe routes through card for recurring
		return []string{"card", "pix"}
	case "boleto":
		return []string{"boleto"}
	default:
		return []string{"card"}
	}
}

func (h *SubscriptionHandler) CreateCheckout(c *gin.Context) {
	userIDVal := utils.GetUserID(c)

	var req createCheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "plano ou tipo de cobrança inválido. Use: pro/premium e monthly/yearly"})
		return
	}

	u, err := h.userRepo.FindByID(userIDVal)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	customerID := u.StripeCustomerID
	if customerID == "" {
		customerID, err = h.gateway.CreateCustomer(u.Email, u.FirstName+" "+u.LastName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create customer"})
			return
		}
		h.db.Model(&user.User{}).Where("id = ?", u.ID).Update("stripe_customer_id", customerID)
	}

	// Determine checkout item id according to payment provider.
	// Stripe => recurring price id
	// AbacatePay => recurring product id (cycle defined in product)
	priceID := h.getCheckoutItemID(req.Plan, req.BillingType)
	if priceID == "" {
		if h.paymentProvider == paymentprovider.AbacatePay {
			c.JSON(http.StatusBadRequest, gin.H{"error": "abacatepay product not configured for this plan. Set ABACATEPAY_PRODUCT_PRO / ABACATEPAY_PRODUCT_PREMIUM (and *_YEARLY)."})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "stripe price not configured for this plan. Set STRIPE_PRICE_PRO or STRIPE_PRICE_PREMIUM env vars."})
		return
	}

	baseURL := os.Getenv("APP_URL")
	result, err := h.gateway.CreateCheckoutSession(payment.CreateCheckoutOptions{
		CustomerID:     customerID,
		PriceID:        priceID,
		SuccessURL:     baseURL + "/subscription?success=true",
		CancelURL:      baseURL + "/subscription?canceled=true",
		Plan:           req.Plan,
		PaymentMethods: resolvePaymentMethods(req.PaymentMethod),
		Metadata:       map[string]string{"user_id": userIDVal.String(), "plan": req.Plan},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create checkout session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"session_id": result.SessionID,
		"url":        result.URL,
	})
}

func (h *SubscriptionHandler) getCheckoutItemID(plan, billingType string) string {
	if h.paymentProvider == paymentprovider.AbacatePay {
		id := h.getAbacateProductIDFromDB(plan, billingType)
		if id != "" {
			return id
		}
		return getAbacateProductIDFromEnv(plan, billingType)
	}
	id := h.getPriceIDFromDB(plan, billingType)
	if id != "" {
		return id
	}
	return getPriceIDFromEnv(plan, billingType)
}

func (h *SubscriptionHandler) getAbacateProductIDFromDB(plan, billingType string) string {
	p, err := h.planRepo.FindBySlug(plan)
	if err != nil {
		return ""
	}
	if billingType == "yearly" {
		return p.StripePriceIDYearly
	}
	return p.StripePriceIDMonthly
}

// getPriceIDFromDB fetches the Stripe Price ID from the plans table.
func (h *SubscriptionHandler) getPriceIDFromDB(plan, billingType string) string {
	p, err := h.planRepo.FindBySlug(plan)
	if err != nil {
		return ""
	}
	if billingType == "yearly" {
		return p.StripePriceIDYearly
	}
	return p.StripePriceIDMonthly
}

// getPriceIDFromEnv falls back to env var configuration (backwards compat).
func getPriceIDFromEnv(plan, billingType string) string {
	envKey := "STRIPE_PRICE_" + strings.ToUpper(plan)
	if billingType == "yearly" {
		envKey += "_YEARLY"
	}
	return os.Getenv(envKey)
}

func getAbacateProductIDFromEnv(plan, billingType string) string {
	envKey := "ABACATEPAY_PRODUCT_" + strings.ToUpper(plan)
	if billingType == "yearly" {
		envKey += "_YEARLY"
	}
	return os.Getenv(envKey)
}

func (h *SubscriptionHandler) CreatePortal(c *gin.Context) {
	if h.paymentProvider == paymentprovider.AbacatePay {
		c.JSON(http.StatusBadRequest, gin.H{"error": "billing portal is unavailable for abacatepay provider"})
		return
	}
	userIDVal := utils.GetUserID(c)

	u, err := h.userRepo.FindByID(userIDVal)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if u.StripeCustomerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no active subscription"})
		return
	}

	baseURL := os.Getenv("APP_URL")
	result, err := h.gateway.CreatePortalSession(u.StripeCustomerID, baseURL+"/subscription")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create portal session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": result.URL})
}

func (h *SubscriptionHandler) HandleStripeWebhook(c *gin.Context) {
	if h.paymentProvider != paymentprovider.Stripe {
		c.JSON(http.StatusBadRequest, gin.H{"error": "stripe webhook disabled for current PAYMENT_PROVIDER"})
		return
	}
	h.handlePaymentWebhook(c, "Stripe-Signature")
}

func (h *SubscriptionHandler) HandleAbacatePayWebhook(c *gin.Context) {
	if h.paymentProvider != paymentprovider.AbacatePay {
		c.JSON(http.StatusBadRequest, gin.H{"error": "abacatepay webhook disabled for current PAYMENT_PROVIDER"})
		return
	}
	h.handlePaymentWebhook(c, "X-Abacate-Webhook-Secret")
}

func (h *SubscriptionHandler) handlePaymentWebhook(c *gin.Context, signatureHeader string) {
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
		return
	}

	signature := c.GetHeader(signatureHeader)
	events, err := h.gateway.ValidateWebhook(payload, signature)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid webhook signature"})
		return
	}

	for _, event := range events {
		if err := h.processWebhookEvent(event); err != nil {
			_ = err
		}
	}

	c.JSON(http.StatusOK, gin.H{"received": true})
}

func (h *SubscriptionHandler) processWebhookEvent(event *payment.WebhookEvent) error {
	switch event.Type {

	case payment.EventCheckoutCompleted:
		// Upsert subscription record with full data from the checkout session
		sub, _ := h.subscriptionRepo.FindByExternalID(event.SubscriptionID)
		if sub == nil {
			sub = &subscription.Subscription{}
		}
		if event.SubscriptionID != "" {
			sub.ExternalID = event.SubscriptionID
		}
		// Resolve Plan Slug from PriceID
		var planSlug string
		planSlug = h.resolvePlanSlug(event.PriceID)

		if planSlug != "" {
			sub.Plan = planSlug
		}
		sub.Status = subscription.StatusActive

		// Set UserID from metadata
		if event.UserID != "" {
			userID, err := uuid.Parse(event.UserID)
			if err == nil {
				sub.UserID = userID
			}
		}

		if err := h.subscriptionRepo.Upsert(sub); err != nil {
			return fmt.Errorf("checkout: upsert subscription: %w", err)
		}

		// Update the user's plan field so future API calls reflect the new plan
		if planSlug != "" {
			if event.UserID != "" {
				h.db.Model(&user.User{}).
					Where("id = ?", event.UserID).
					Update("plan", planSlug)
			} else if event.CustomerID != "" && h.paymentProvider == paymentprovider.Stripe {
				h.db.Model(&user.User{}).
					Where("stripe_customer_id = ?", event.CustomerID).
					Update("plan", planSlug)
			}
		}

	case payment.EventSubscriptionUpdated:
		sub, _ := h.subscriptionRepo.FindByExternalID(event.SubscriptionID)
		if sub == nil {
			sub = &subscription.Subscription{}
		}
		sub.ExternalID = event.SubscriptionID
		var planSlug string
		planSlug = h.resolvePlanSlug(event.PriceID)
		if planSlug != "" {
			sub.Plan = planSlug
		}
		sub.Status = subscription.SubscriptionStatus(event.Status)

		if err := h.subscriptionRepo.Upsert(sub); err != nil {
			return err
		}
		if planSlug != "" && event.CustomerID != "" && h.paymentProvider == paymentprovider.Stripe {
			h.db.Model(&user.User{}).
				Where("stripe_customer_id = ?", event.CustomerID).
				Update("plan", planSlug)
		}

	case payment.EventSubscriptionDeleted:
		sub, err := h.subscriptionRepo.FindByExternalID(event.SubscriptionID)
		if err == nil {
			sub.Status = subscription.StatusCanceled
			_ = h.subscriptionRepo.Upsert(sub)
		}
		if h.paymentProvider == paymentprovider.Stripe {
			h.db.Model(&user.User{}).
				Where("stripe_customer_id = ?", event.CustomerID).
				Update("plan", "free")
		} else if event.UserID != "" {
			h.db.Model(&user.User{}).
				Where("id = ?", event.UserID).
				Update("plan", "free")
		}
	}

	return nil
}

func (h *SubscriptionHandler) resolvePlanSlug(priceOrSlug string) string {
	if priceOrSlug == "" {
		return ""
	}
	if priceOrSlug == "pro" || priceOrSlug == "premium" || priceOrSlug == "free" {
		return priceOrSlug
	}
	if p, err := h.planRepo.FindByStripePriceID(priceOrSlug); err == nil && p != nil {
		return p.Slug
	}
	return ""
}
