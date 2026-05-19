package stripe

import (
	"encoding/json"
	"finance-ia/internal/domain/payment"
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/stripe/stripe-go/v76"
	portalsession "github.com/stripe/stripe-go/v76/billingportal/session"
	checkoutsession "github.com/stripe/stripe-go/v76/checkout/session"
	"github.com/stripe/stripe-go/v76/customer"
	"github.com/stripe/stripe-go/v76/invoice"
	"github.com/stripe/stripe-go/v76/price"
	"github.com/stripe/stripe-go/v76/product"
	stripeSubscription "github.com/stripe/stripe-go/v76/subscription"
	"github.com/stripe/stripe-go/v76/webhook"
)

// Adapter implements the payment.PaymentGateway port for Stripe.
// To replace Stripe with another provider, implement payment.PaymentGateway without touching this file.
type Adapter struct {
	webhookSecret string
}

// NewAdapter creates a new Stripe payment gateway adapter
func NewAdapter() *Adapter {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	return &Adapter{
		webhookSecret: os.Getenv("STRIPE_WEBHOOK_SECRET"),
	}
}

func (a *Adapter) CreateCheckoutSession(opts payment.CreateCheckoutOptions) (*payment.CheckoutResult, error) {
	params := &stripe.CheckoutSessionParams{
		Customer:   stripe.String(opts.CustomerID),
		Mode:       stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		SuccessURL: stripe.String(opts.SuccessURL),
		CancelURL:  stripe.String(opts.CancelURL),
		// Not setting PaymentMethodTypes — Stripe will automatically show all payment
		// methods enabled in the Dashboard (card, PIX, Boleto, etc.)
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(opts.PriceID),
				Quantity: stripe.Int64(1),
			},
		},
	}

	// Embed user + plan metadata for webhook identification
	if len(opts.Metadata) > 0 {
		params.Metadata = make(map[string]string)
		for k, v := range opts.Metadata {
			params.Metadata[k] = v
		}
	}

	s, err := checkoutsession.New(params)
	if err != nil {
		return nil, fmt.Errorf("stripe: create checkout session: %w", err)
	}

	return &payment.CheckoutResult{
		SessionID: s.ID,
		URL:       s.URL,
	}, nil
}

func (a *Adapter) CreatePortalSession(customerID string, returnURL string) (*payment.PortalResult, error) {
	params := &stripe.BillingPortalSessionParams{
		Customer:  stripe.String(customerID),
		ReturnURL: stripe.String(returnURL),
	}

	s, err := portalsession.New(params)
	if err != nil {
		return nil, fmt.Errorf("stripe: create portal session: %w", err)
	}

	return &payment.PortalResult{URL: s.URL}, nil
}

func (a *Adapter) CreateCustomer(email, name string) (string, error) {
	params := &stripe.CustomerParams{
		Email: stripe.String(email),
		Name:  stripe.String(name),
	}

	c, err := customer.New(params)
	if err != nil {
		return "", fmt.Errorf("stripe: create customer: %w", err)
	}

	return c.ID, nil
}

func (a *Adapter) GetSubscription(externalID string) (*payment.SubscriptionInfo, error) {
	sub, err := stripeSubscription.Get(externalID, nil)
	if err != nil {
		return nil, fmt.Errorf("stripe: get subscription: %w", err)
	}

	priceID := ""
	if len(sub.Items.Data) > 0 {
		priceID = sub.Items.Data[0].Price.ID
	}

	return &payment.SubscriptionInfo{
		ExternalID:  sub.ID,
		CustomerID:  sub.Customer.ID,
		PriceID:     priceID,
		Status:      string(sub.Status),
		PeriodStart: sub.CurrentPeriodStart,
		PeriodEnd:   sub.CurrentPeriodEnd,
	}, nil
}

func (a *Adapter) CancelSubscription(externalID string) error {
	params := &stripe.SubscriptionParams{
		CancelAtPeriodEnd: stripe.Bool(true),
	}
	_, err := stripeSubscription.Update(externalID, params)
	if err != nil {
		return fmt.Errorf("stripe: cancel subscription: %w", err)
	}
	return nil
}

func (a *Adapter) CreateProduct(name, description string) (string, error) {
	params := &stripe.ProductParams{
		Name:        stripe.String(name),
		Description: stripe.String(description),
	}
	p, err := product.New(params)
	if err != nil {
		return "", fmt.Errorf("stripe: create product: %w", err)
	}
	return p.ID, nil
}

func (a *Adapter) CreatePrice(productID string, amount float64, currency string, interval string) (string, error) {
	params := &stripe.PriceParams{
		Product:    stripe.String(productID),
		UnitAmount: stripe.Int64(int64(amount * 100)), // Stripe uses cents
		Currency:   stripe.String(currency),
		Recurring: &stripe.PriceRecurringParams{
			Interval: stripe.String(interval), // "month" or "year"
		},
	}
	p, err := price.New(params)
	if err != nil {
		return "", fmt.Errorf("stripe: create price: %w", err)
	}
	return p.ID, nil
}

func (a *Adapter) ArchiveProduct(productID string) error {
	params := &stripe.ProductParams{
		Active: stripe.Bool(false),
	}
	_, err := product.Update(productID, params)
	if err != nil {
		return fmt.Errorf("stripe: archive product: %w", err)
	}
	return nil
}

func (a *Adapter) ValidateWebhook(payload []byte, signature string) ([]*payment.WebhookEvent, error) {
	event, err := webhook.ConstructEvent(payload, signature, a.webhookSecret)
	if err != nil {
		return nil, fmt.Errorf("stripe: webhook validation failed: %w", err)
	}

	var events []*payment.WebhookEvent

	switch event.Type {
	case "checkout.session.completed":
		var s stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &s); err != nil {
			return nil, err
		}
		customerID := ""
		if s.Customer != nil {
			customerID = s.Customer.ID
		}
		subscriptionID := ""
		if s.Subscription != nil {
			subscriptionID = s.Subscription.ID
		}
		userID := ""
		plan := ""
		if s.Metadata != nil {
			userID = s.Metadata["user_id"]
			plan = s.Metadata["plan"]
		}
		events = append(events, &payment.WebhookEvent{
			Type:           payment.EventCheckoutCompleted,
			CustomerID:     customerID,
			SubscriptionID: subscriptionID,
			UserID:         userID,
			PriceID:        plan, // We stored plan slug here previously in metadata, we'll keep as is for backward compat or update later. Since we now use price IDs primarily, checkout doesn't map directly from metadata but from the recurring item if needed. For now we will map what's in metadata as PriceID for compat or we can fetch the price. Checkout metadata usually had 'plan' slug.
			Status:         "active",
		})

	case "customer.subscription.updated":
		var sub stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &sub); err != nil {
			return nil, err
		}
		priceID := ""
		if len(sub.Items.Data) > 0 {
			priceID = sub.Items.Data[0].Price.ID
		}
		customerID := ""
		if sub.Customer != nil {
			customerID = sub.Customer.ID
		}
		events = append(events, &payment.WebhookEvent{
			Type:           payment.EventSubscriptionUpdated,
			SubscriptionID: sub.ID,
			CustomerID:     customerID,
			PriceID:        priceID,
			Status:         string(sub.Status),
		})

	case "customer.subscription.deleted":
		var sub stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &sub); err != nil {
			return nil, err
		}
		customerID := ""
		if sub.Customer != nil {
			customerID = sub.Customer.ID
		}
		events = append(events, &payment.WebhookEvent{
			Type:           payment.EventSubscriptionDeleted,
			SubscriptionID: sub.ID,
			CustomerID:     customerID,
			Status:         "canceled",
		})
	}

	return events, nil
}

func (a *Adapter) ListInvoices(customerID string) ([]*payment.InvoiceInfo, error) {
	params := &stripe.InvoiceListParams{
		Customer: stripe.String(customerID),
	}
	params.Filters.AddFilter("limit", "", "10")

	iter := invoice.List(params)
	var result []*payment.InvoiceInfo
	for iter.Next() {
		inv := iter.Invoice()
		result = append(result, &payment.InvoiceInfo{
			ID:               inv.ID,
			Date:             time.Unix(inv.Created, 0).Format("02/01/2006"),
			Amount:           float64(inv.AmountPaid) / 100.0,
			Status:           string(inv.Status),
			Currency:         string(inv.Currency),
			HostedInvoiceURL: inv.HostedInvoiceURL,
		})
	}
	if err := iter.Err(); err != nil {
		return nil, fmt.Errorf("stripe: list invoices: %w", err)
	}
	return result, nil
}

// ListActiveSubscriptions lists active (or trialing) subscriptions for a customer.
// Used to sync user's plan when webhooks were missed.
func (a *Adapter) ListActiveSubscriptions(customerID string) ([]*payment.SubscriptionInfo, error) {
	params := &stripe.SubscriptionListParams{
		Customer: stripe.String(customerID),
		Status:   stripe.String("active"),
	}

	iter := stripeSubscription.List(params)
	var result []*payment.SubscriptionInfo
	for iter.Next() {
		sub := iter.Subscription()
		priceID := ""
		if len(sub.Items.Data) > 0 {
			priceID = sub.Items.Data[0].Price.ID
		}
		customID := ""
		if sub.Customer != nil {
			customID = sub.Customer.ID
		}
		result = append(result, &payment.SubscriptionInfo{
			ExternalID:  sub.ID,
			CustomerID:  customID,
			PriceID:     priceID,
			Status:      string(sub.Status),
			PeriodStart: sub.CurrentPeriodStart,
			PeriodEnd:   sub.CurrentPeriodEnd,
		})
	}
	if err := iter.Err(); err != nil {
		return nil, fmt.Errorf("stripe: list active subscriptions: %w", err)
	}
	return result, nil
}

const planSlugMetadataKey = "finance_plan_slug"

// EnsurePlanCatalog guarantees that a paid plan has one product and two recurring prices in Stripe.
// It is idempotent: finds reusable objects first, creates only when missing.
func (a *Adapter) EnsurePlanCatalog(slug, name, description string, monthly, yearly float64, currency string) (string, string, string, error) {
	productID, err := a.findProductByPlanSlug(slug)
	if err != nil {
		return "", "", "", err
	}
	if productID == "" {
		params := &stripe.ProductParams{
			Name:        stripe.String(name),
			Description: stripe.String(description),
			Metadata: map[string]string{
				planSlugMetadataKey: slug,
			},
		}
		p, createErr := product.New(params)
		if createErr != nil {
			return "", "", "", fmt.Errorf("stripe: create product: %w", createErr)
		}
		productID = p.ID
	}

	monthlyID, err := a.findRecurringPriceID(productID, monthly, currency, "month")
	if err != nil {
		return "", "", "", err
	}
	if monthlyID == "" {
		monthlyID, err = a.CreatePrice(productID, monthly, currency, "month")
		if err != nil {
			return "", "", "", err
		}
	}

	yearlyID, err := a.findRecurringPriceID(productID, yearly, currency, "year")
	if err != nil {
		return "", "", "", err
	}
	if yearlyID == "" {
		yearlyID, err = a.CreatePrice(productID, yearly, currency, "year")
		if err != nil {
			return "", "", "", err
		}
	}

	return productID, monthlyID, yearlyID, nil
}

func (a *Adapter) findProductByPlanSlug(slug string) (string, error) {
	params := &stripe.ProductListParams{}
	params.Filters.AddFilter("limit", "", "100")
	iter := product.List(params)
	for iter.Next() {
		p := iter.Product()
		if p != nil && p.Metadata != nil && p.Metadata[planSlugMetadataKey] == slug {
			return p.ID, nil
		}
	}
	if err := iter.Err(); err != nil {
		return "", fmt.Errorf("stripe: list products: %w", err)
	}
	return "", nil
}

func (a *Adapter) findRecurringPriceID(productID string, amount float64, currency string, interval string) (string, error) {
	params := &stripe.PriceListParams{
		Product: stripe.String(productID),
		Active:  stripe.Bool(true),
	}
	params.Filters.AddFilter("limit", "", "100")

	wantCents := int64(math.Round(amount * 100))
	wantCurrency := strings.ToLower(currency)
	iter := price.List(params)
	for iter.Next() {
		p := iter.Price()
		if p == nil || p.Recurring == nil {
			continue
		}
		if p.UnitAmount == wantCents &&
			strings.EqualFold(string(p.Currency), wantCurrency) &&
			string(p.Recurring.Interval) == interval {
			return p.ID, nil
		}
	}
	if err := iter.Err(); err != nil {
		return "", fmt.Errorf("stripe: list prices: %w", err)
	}
	return "", nil
}
