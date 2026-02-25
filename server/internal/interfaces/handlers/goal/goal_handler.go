package goal

import (
	"finance-ia/internal/domain/goal"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GoalHandler struct {
	service *goal.Service
}

func NewGoalHandler(service *goal.Service) *GoalHandler {
	return &GoalHandler{service: service}
}

func (h *GoalHandler) RegisterRoutes(public, protected gin.IRouter) {
	g := protected.Group("/goals")
	{
		g.POST("/", h.CreateGoal)
		g.GET("/", h.GetGoals)
		g.PUT("/:id", h.UpdateGoal)
		g.DELETE("/:id", h.DeleteGoal)
	}
}

// requireProOrPremium was removed to allow free users access up to limits

func (h *GoalHandler) CreateGoal(c *gin.Context) {
	var input struct {
		Name         string  `json:"name" binding:"required"`
		TargetAmount float64 `json:"target_amount" binding:"required"`
		TargetDate   string  `json:"target_date" binding:"required"`
		Icon         string  `json:"icon"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := h.getUserID(c)

	targetDate, err := http.ParseTime(input.TargetDate)
	if err != nil {
		// fallback simple parse
		targetDate, _ = http.ParseTime(input.TargetDate + "T00:00:00Z")
	}

	// Validate Free Tier Limit
	planRaw, _ := c.Get("plan")
	if planStr, ok := planRaw.(string); ok && planStr == "free" || planStr == "" {
		existingGoals, err := h.service.GetByUser(userID)
		if err == nil && len(existingGoals) >= 2 {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "O plano gratuito permite criar no máximo 2 metas. Faça upgrade para o Pro para metas ilimitadas.",
				"upgrade": true,
			})
			return
		}
	}

	newGoal := &goal.Goal{
		UserID:       userID,
		Name:         input.Name,
		TargetAmount: input.TargetAmount,
		TargetDate:   targetDate,
		Icon:         input.Icon,
	}
	if newGoal.Icon == "" {
		newGoal.Icon = "flag"
	}

	if err := h.service.Create(newGoal); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newGoal)
}

func (h *GoalHandler) GetGoals(c *gin.Context) {
	userID := h.getUserID(c)
	goals, err := h.service.GetByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if goals == nil {
		goals = []*goal.Goal{}
	}
	c.JSON(http.StatusOK, goals)
}

func (h *GoalHandler) UpdateGoal(c *gin.Context) {
	idStr := c.Param("id")
	goalID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid goal ID"})
		return
	}

	var input struct {
		Name          string  `json:"name"`
		TargetAmount  float64 `json:"target_amount"`
		CurrentAmount float64 `json:"current_amount"`
		TargetDate    string  `json:"target_date"`
		Icon          string  `json:"icon"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingGoal, err := h.service.GetByID(goalID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Goal not found"})
		return
	}

	userID := h.getUserID(c)
	if existingGoal.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if input.Name != "" {
		existingGoal.Name = input.Name
	}
	if input.TargetAmount > 0 {
		existingGoal.TargetAmount = input.TargetAmount
	}
	// Can be 0
	existingGoal.CurrentAmount = input.CurrentAmount

	if input.TargetDate != "" {
		targetDate, _ := http.ParseTime(input.TargetDate)
		existingGoal.TargetDate = targetDate
	}
	if input.Icon != "" {
		existingGoal.Icon = input.Icon
	}

	if err := h.service.Update(existingGoal); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, existingGoal)
}

func (h *GoalHandler) DeleteGoal(c *gin.Context) {
	idStr := c.Param("id")
	goalID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid goal ID"})
		return
	}

	existingGoal, err := h.service.GetByID(goalID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Goal not found"})
		return
	}

	userID := h.getUserID(c)
	if existingGoal.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if err := h.service.Delete(goalID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Goal deleted successfully"})
}

func (h *GoalHandler) getUserID(c *gin.Context) uuid.UUID {
	userIDstr, _ := c.Get("user_id")
	id, _ := uuid.Parse(userIDstr.(string))
	return id
}
