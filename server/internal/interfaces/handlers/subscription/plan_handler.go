package subscription

import (
	"finance-ia/internal/domain/payment"
	"finance-ia/internal/domain/subscription"
	dbsub "finance-ia/internal/infrastructure/database/subscription"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PlanHandler struct {
	gateway  payment.PaymentGateway
	planRepo *dbsub.PlanRepository
}

func NewPlanHandler(gateway payment.PaymentGateway, planRepo *dbsub.PlanRepository) *PlanHandler {
	return &PlanHandler{
		gateway:  gateway,
		planRepo: planRepo,
	}
}

func (h *PlanHandler) RegisterRoutes(public gin.IRouter, protected gin.IRouter) {
	plans := protected.Group("/admin/plans")
	{
		plans.POST("/", h.CreatePlan)
		plans.GET("/", h.ListPlans)
		plans.GET("/:id", h.GetPlan)
		plans.PUT("/:id", h.UpdatePlan)
		plans.DELETE("/:id", h.DeletePlan)

		// Feature specific endpoints
		plans.POST("/:id/features", h.AddFeature)
		plans.DELETE("/features/:featureId", h.RemoveFeature)
	}
}

type createPlanRequest struct {
	Slug         string   `json:"slug" binding:"required"`
	Name         string   `json:"name" binding:"required"`
	Description  string   `json:"description"`
	PriceMonthly float64  `json:"price_monthly"`
	PriceYearly  float64  `json:"price_yearly"`
	Features     []string `json:"features"`
}

func (h *PlanHandler) CreatePlan(c *gin.Context) {
	var req createPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. Create Product in Stripe
	stripeProductID, err := h.gateway.CreateProduct(req.Name, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create product in stripe: " + err.Error()})
		return
	}

	// 2. Create Prices in Stripe
	var stripePriceMonthly, stripePriceYearly string
	if req.PriceMonthly > 0 {
		stripePriceMonthly, _ = h.gateway.CreatePrice(stripeProductID, req.PriceMonthly, "brl", "month")
	}
	if req.PriceYearly > 0 {
		stripePriceYearly, _ = h.gateway.CreatePrice(stripeProductID, req.PriceYearly, "brl", "year")
	}

	// 3. Prepare features
	features := make([]subscription.PlanFeature, len(req.Features))
	for i, f := range req.Features {
		features[i] = subscription.PlanFeature{Description: f}
	}

	// 4. Save in DB
	plan := &subscription.Plan{
		Slug:                 req.Slug,
		Name:                 req.Name,
		Description:          req.Description,
		PriceMonthly:         req.PriceMonthly,
		PriceYearly:          req.PriceYearly,
		StripeProductID:      stripeProductID,
		StripePriceIDMonthly: stripePriceMonthly,
		StripePriceIDYearly:  stripePriceYearly,
		Features:             features,
		IsActive:             true,
	}

	if err := h.planRepo.Upsert(plan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save plan in database"})
		return
	}

	c.JSON(http.StatusCreated, plan)
}

func (h *PlanHandler) ListPlans(c *gin.Context) {
	plans, err := h.planRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch plans"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"plans": plans})
}

func (h *PlanHandler) GetPlan(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	plan, err := h.planRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "plan not found"})
		return
	}

	c.JSON(http.StatusOK, plan)
}

func (h *PlanHandler) UpdatePlan(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req createPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plan, err := h.planRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "plan not found"})
		return
	}

	// If prices have changed, we need to create new Prices in Stripe
	if req.PriceMonthly != plan.PriceMonthly {
		if req.PriceMonthly > 0 {
			newPriceID, err := h.gateway.CreatePrice(plan.StripeProductID, req.PriceMonthly, "brl", "month")
			if err == nil {
				plan.StripePriceIDMonthly = newPriceID
			}
		}
		plan.PriceMonthly = req.PriceMonthly
	}

	if req.PriceYearly != plan.PriceYearly {
		if req.PriceYearly > 0 {
			newPriceID, err := h.gateway.CreatePrice(plan.StripeProductID, req.PriceYearly, "brl", "year")
			if err == nil {
				plan.StripePriceIDYearly = newPriceID
			}
		}
		plan.PriceYearly = req.PriceYearly
	}

	plan.Slug = req.Slug
	plan.Name = req.Name
	plan.Description = req.Description

	if err := h.planRepo.Upsert(plan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update plan"})
		return
	}

	c.JSON(http.StatusOK, plan)
}

func (h *PlanHandler) DeletePlan(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.planRepo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete plan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "plan deleted"})
}

type addFeatureRequest struct {
	Description string `json:"description" binding:"required"`
}

func (h *PlanHandler) AddFeature(c *gin.Context) {
	planID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan id"})
		return
	}

	var req addFeatureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	feature := &subscription.PlanFeature{
		PlanID:      planID,
		Description: req.Description,
	}

	if err := h.planRepo.UpsertFeature(feature); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add feature"})
		return
	}

	c.JSON(http.StatusCreated, feature)
}

func (h *PlanHandler) RemoveFeature(c *gin.Context) {
	featureID, err := uuid.Parse(c.Param("featureId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid feature id"})
		return
	}

	if err := h.planRepo.DeleteFeature(featureID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove feature"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "feature removed"})
}
