package finance

import (
	domainFinance "finance-ia/internal/domain/finance"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// BehaviorHandler provides advanced behavioral analysis (premium+)
type BehaviorHandler struct {
	financeService *domainFinance.Service
}

func NewBehaviorHandler(svc *domainFinance.Service) *BehaviorHandler {
	return &BehaviorHandler{financeService: svc}
}

func (h *BehaviorHandler) RegisterRoutes(public, protected gin.IRouter) {
	g := protected.Group("/finance")
	g.GET("/behavior", h.GetBehavior)
}

// GetBehavior returns the behavioral analysis for the authenticated user.
// Requires "premium" or "pro" plan — returns 403 for free users.
func (h *BehaviorHandler) GetBehavior(c *gin.Context) {
	plan, _ := c.Get("plan")
	planStr, _ := plan.(string)

	if planStr == "free" || planStr == "" {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "upgrade_required: behavioral analysis requires Premium or Pro plan",
			"upgrade": true,
		})
		return
	}

	userIDStr, _ := c.Get("user_id")
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	analysis, err := h.financeService.AnalyzeBehavior(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, analysis)
}
