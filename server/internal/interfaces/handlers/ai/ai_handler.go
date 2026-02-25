package ai

import (
	"finance-ia/internal/domain/ai"
	"finance-ia/internal/domain/finance"
	"finance-ia/internal/domain/user"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AIHandler struct {
	aiService      *ai.Service
	financeService *finance.Service
	userService    *user.Service
}

func NewAIHandler(aiService *ai.Service, financeService *finance.Service, userService *user.Service) *AIHandler {
	return &AIHandler{aiService: aiService, financeService: financeService, userService: userService}
}

func (h *AIHandler) RegisterRoutes(public, protected gin.IRouter) {
	g := protected.Group("/ai")
	{
		g.GET("/insight", h.GetInsight)
		g.GET("/diagnostic", h.GetDiagnostic)
		g.GET("/simulator", h.GetSimulator)
		g.GET("/coach", h.GetCoachStatus)
		g.POST("/coach", h.PostCoachMessage)
		g.GET("/missions", h.GetMissions)
		g.GET("/health-score", h.GetHealthScore)
	}
}

func (h *AIHandler) GetInsight(c *gin.Context) {
	userID := getUserID(c)

	// Buscar o plano real do BD invés de confiar no JWT
	u, err := h.userService.GetByID(userID)
	plan := "free"
	if err == nil && u != nil {
		plan = string(u.Plan)
	}

	// Build context from real financial data
	ctx, err := h.buildFinancialContext(userID, plan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load financial data"})
		return
	}

	insight, err := h.aiService.GetInsight(userID, plan, *ctx)
	if err != nil {
		if err.Error()[:12] == "rate_limited" {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   err.Error(),
				"upgrade": true,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, insight)
}

func (h *AIHandler) GetDiagnostic(c *gin.Context) {
	userID := getUserID(c)

	u, err := h.userService.GetByID(userID)
	plan := "free"
	if err == nil && u != nil {
		plan = string(u.Plan)
	}

	ctx, err := h.buildFinancialContext(userID, plan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load financial data"})
		return
	}

	insight, err := h.aiService.GetDiagnostic(userID, plan, *ctx)
	if err != nil {
		if err.Error()[:17] == "upgrade_required:" {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   err.Error(),
				"upgrade": true,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"diagnostic": insight})
}

func (h *AIHandler) GetSimulator(c *gin.Context) {
	userID := getUserID(c)

	u, err := h.userService.GetByID(userID)
	plan := "free"
	if err == nil && u != nil {
		plan = string(u.Plan)
	}

	ctx, err := h.buildFinancialContext(userID, plan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load financial data"})
		return
	}

	insight, err := h.aiService.GetProjection(userID, plan, *ctx)
	if err != nil {
		if err.Error()[:17] == "upgrade_required:" {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   err.Error(),
				"upgrade": true,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"projection": insight})
}

func (h *AIHandler) GetCoachStatus(c *gin.Context) {
	userID := getUserID(c)

	u, err := h.userService.GetByID(userID)
	if err != nil || u == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	maxUses := 50
	now := time.Now()

	if now.Sub(u.LastAICoachReset) > 7*24*time.Hour {
		u.AICoachUsesThisWeek = 0
		u.LastAICoachReset = now
		h.userService.Update(u) // Ignora erro de persistência silenciosa
	}

	remaining := maxUses - u.AICoachUsesThisWeek
	if remaining < 0 {
		remaining = 0
	}

	c.JSON(http.StatusOK, gin.H{
		"remaining": remaining,
		"max":       maxUses,
		"uses":      u.AICoachUsesThisWeek,
	})
}

func (h *AIHandler) PostCoachMessage(c *gin.Context) {
	userID := getUserID(c)

	u, err := h.userService.GetByID(userID)
	if err != nil || u == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	if u.Plan != "premium" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Upgrade required: Coach is Premium only.", "upgrade": true})
		return
	}

	maxUses := 50
	now := time.Now()

	// Check weekly reset
	if now.Sub(u.LastAICoachReset) > 7*24*time.Hour {
		u.AICoachUsesThisWeek = 0
		u.LastAICoachReset = now
	}

	if u.AICoachUsesThisWeek >= maxUses {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error":     "Limite semanal do Coach atingido. Tente novamente próxima semana.",
			"remaining": 0,
			"max":       maxUses,
		})
		return
	}

	// Increment usages
	u.AICoachUsesThisWeek++
	if err := h.userService.Update(u); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to secure AI token"})
		return
	}

	remaining := maxUses - u.AICoachUsesThisWeek

	// Here would go the actual logic to call the AI provider to chat
	// that we will implement next.

	c.JSON(http.StatusOK, gin.H{
		"message":   "Sua mensagem foi recebida pela IA (Em Desenvolvimento).",
		"remaining": remaining,
		"max":       maxUses,
	})
}

func (h *AIHandler) GetMissions(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Missões em desenvolvimento."})
}

func (h *AIHandler) GetHealthScore(c *gin.Context) {
	userID := getUserID(c)

	summary, err := h.financeService.GetDashboardSummary(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"score":        summary.HealthScore,
		"level":        summary.HealthLevel,
		"savings_rate": summary.SavingsRate,
	})
}

func (h *AIHandler) buildFinancialContext(userID uuid.UUID, plan string) (*ai.FinancialContext, error) {
	summary, err := h.financeService.GetDashboardSummary(userID)
	if err != nil {
		return nil, err
	}

	ctx := &ai.FinancialContext{
		Plan:          plan,
		TotalIncome:   summary.TotalIncome,
		TotalExpenses: summary.TotalExpenses,
		Balance:       summary.Balance,
		SavingsRate:   summary.SavingsRate,
		HealthScore:   summary.HealthScore,
		HealthLevel:   summary.HealthLevel,
	}

	for _, cs := range summary.CategoryBreakdown {
		ctx.TopCategories = append(ctx.TopCategories, ai.CategorySpend{
			Name:       cs.CategoryName,
			Amount:     cs.Total,
			Percentage: cs.Percentage,
		})
	}

	for _, mt := range summary.MonthlyTrend {
		ctx.MonthlyTrends = append(ctx.MonthlyTrends, ai.MonthTrend{
			Month:    mt.Month,
			Income:   mt.Income,
			Expenses: mt.Expenses,
		})
	}

	return ctx, nil
}

func getUserID(c *gin.Context) uuid.UUID {
	userIDstr, _ := c.Get("user_id")
	id, _ := uuid.Parse(userIDstr.(string))
	return id
}
