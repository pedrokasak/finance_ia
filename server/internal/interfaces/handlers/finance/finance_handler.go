package finance

import (
	"finance-ia/internal/domain/finance"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FinanceHandler struct {
	financeService *finance.Service
}

func NewFinanceHandler(financeService *finance.Service) *FinanceHandler {
	return &FinanceHandler{financeService: financeService}
}

func (h *FinanceHandler) RegisterRoutes(public, protected gin.IRouter) {
	g := protected.Group("/finance")
	{
		g.POST("/transactions", h.CreateTransaction)
		g.GET("/transactions", h.ListTransactions)
		g.PUT("/transactions/:id", h.UpdateTransaction)
		g.DELETE("/transactions/:id", h.DeleteTransaction)
		g.GET("/categories", h.ListCategories)
		g.POST("/categories", h.CreateCategory)
		g.PUT("/categories/:id", h.UpdateCategory)
		g.DELETE("/categories/:id", h.DeleteCategory)
		g.GET("/budget", h.GetBudget)
		g.POST("/budget", h.UpsertBudget)
		g.GET("/dashboard", h.GetDashboard)
		g.GET("/methods", h.ListFinancialMethods)
	}
}

type createTransactionRequest struct {
	CategoryID     *string `json:"category_id"`
	Type           string  `json:"type" binding:"required,oneof=income expense"`
	Amount         float64 `json:"amount" binding:"required,gt=0"`
	Description    string  `json:"description"`
	Date           string  `json:"date"` // ISO 8601
	IsRecurring    bool    `json:"is_recurring"`
	IdempotencyKey string  `json:"idempotency_key"`
}

func (h *FinanceHandler) CreateTransaction(c *gin.Context) {
	userID := getUserID(c)

	var req createTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := &finance.Transaction{
		UserID:         userID,
		Type:           finance.TransactionType(req.Type),
		Amount:         req.Amount,
		Description:    req.Description,
		IsRecurring:    req.IsRecurring,
		IdempotencyKey: req.IdempotencyKey,
		Date:           time.Now(),
	}

	if req.Date != "" {
		if d, err := time.Parse("2006-01-02", req.Date); err == nil {
			tx.Date = d
		}
	}

	if req.CategoryID != nil {
		if catID, err := uuid.Parse(*req.CategoryID); err == nil {
			tx.CategoryID = &catID
		}
	}

	// Also check Idempotency-Key header (middleware populates it)
	if headerKey := c.GetHeader("Idempotency-Key"); headerKey != "" && tx.IdempotencyKey == "" {
		tx.IdempotencyKey = headerKey
	}

	if err := h.financeService.CreateTransaction(tx); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tx)
}

type updateTransactionRequest struct {
	CategoryID  *string `json:"category_id"`
	Type        string  `json:"type" binding:"required,oneof=income expense"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Description string  `json:"description"`
	Date        string  `json:"date"` // ISO 8601
	IsRecurring bool    `json:"is_recurring"`
}

func (h *FinanceHandler) UpdateTransaction(c *gin.Context) {
	userID := getUserID(c)
	txID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction id"})
		return
	}

	var req updateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payload := &finance.Transaction{
		Type:        finance.TransactionType(req.Type),
		Amount:      req.Amount,
		Description: req.Description,
		IsRecurring: req.IsRecurring,
	}

	if req.Date != "" {
		if d, err := time.Parse("2006-01-02", req.Date); err == nil {
			payload.Date = d
		}
	}

	if req.CategoryID != nil {
		if catID, err := uuid.Parse(*req.CategoryID); err == nil {
			payload.CategoryID = &catID
		}
	}

	updatedTx, err := h.financeService.UpdateTransaction(txID, userID, payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedTx)
}

func (h *FinanceHandler) ListTransactions(c *gin.Context) {
	userID := getUserID(c)

	filter := finance.TransactionFilter{
		UserID: userID,
		Page:   1,
		Limit:  20,
	}

	if p, err := strconv.Atoi(c.Query("page")); err == nil && p > 0 {
		filter.Page = p
	}
	if l, err := strconv.Atoi(c.Query("limit")); err == nil && l > 0 {
		filter.Limit = l
	}
	if t := c.Query("type"); t != "" {
		txType := finance.TransactionType(t)
		filter.Type = &txType
	}
	if catStr := c.Query("category_id"); catStr != "" {
		if catID, err := uuid.Parse(catStr); err == nil {
			filter.CategoryID = &catID
		}
	}
	if s := c.Query("start_date"); s != "" {
		if d, err := time.Parse("2006-01-02", s); err == nil {
			filter.StartDate = &d
		}
	}
	if e := c.Query("end_date"); e != "" {
		if d, err := time.Parse("2006-01-02", e); err == nil {
			filter.EndDate = &d
		}
	}

	txs, total, err := h.financeService.ListTransactions(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  txs,
		"total": total,
		"page":  filter.Page,
		"limit": filter.Limit,
	})
}

func (h *FinanceHandler) DeleteTransaction(c *gin.Context) {
	userID := getUserID(c)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.financeService.DeleteTransaction(id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "transaction deleted"})
}

func (h *FinanceHandler) ListCategories(c *gin.Context) {
	userID := getUserID(c)
	cats, err := h.financeService.GetCategories(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cats)
}

func (h *FinanceHandler) ListFinancialMethods(c *gin.Context) {
	methods, err := h.financeService.GetFinancialMethods()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, methods)
}

type createCategoryRequest struct {
	Name  string `json:"name" binding:"required"`
	Type  string `json:"type" binding:"required,oneof=income expense"`
	Color string `json:"color"`
	Icon  string `json:"icon"`
}

func (h *FinanceHandler) CreateCategory(c *gin.Context) {
	userID := getUserID(c)
	var req createCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cat := &finance.Category{
		UserID:    &userID,
		Name:      req.Name,
		Type:      finance.CategoryType(req.Type),
		Color:     req.Color,
		Icon:      req.Icon,
		IsDefault: false,
	}

	if err := h.financeService.CreateCategory(cat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, cat)
}

type updateCategoryRequest struct {
	Name  string `json:"name"`
	Color string `json:"color"`
	Icon  string `json:"icon"`
}

func (h *FinanceHandler) UpdateCategory(c *gin.Context) {
	userID := getUserID(c)
	catID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
		return
	}

	var req updateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payload := &finance.Category{
		Name:  req.Name,
		Color: req.Color,
		Icon:  req.Icon,
	}

	updatedCat, err := h.financeService.UpdateCategory(catID, userID, payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedCat)
}

func (h *FinanceHandler) DeleteCategory(c *gin.Context) {
	userID := getUserID(c)
	catID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
		return
	}

	if err := h.financeService.DeleteCategory(catID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "category deleted"})
}

func (h *FinanceHandler) GetBudget(c *gin.Context) {
	userID := getUserID(c)
	budget, err := h.financeService.GetCurrentBudget(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "no budget set for this period"})
		return
	}
	c.JSON(http.StatusOK, budget)
}

type upsertBudgetRequest struct {
	TotalIncome    float64 `json:"total_income" binding:"required,gt=0"`
	NeedsPercent   float64 `json:"needs_percent"`
	WantsPercent   float64 `json:"wants_percent"`
	SavingsPercent float64 `json:"savings_percent"`
}

func (h *FinanceHandler) UpsertBudget(c *gin.Context) {
	userID := getUserID(c)
	var req upsertBudgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Default to 50-30-20 if not specified
	if req.NeedsPercent == 0 && req.WantsPercent == 0 && req.SavingsPercent == 0 {
		req.NeedsPercent = 50
		req.WantsPercent = 30
		req.SavingsPercent = 20
	}

	budget := &finance.Budget{
		UserID:         userID,
		Period:         time.Now().Format("2006-01"),
		TotalIncome:    req.TotalIncome,
		NeedsPercent:   req.NeedsPercent,
		WantsPercent:   req.WantsPercent,
		SavingsPercent: req.SavingsPercent,
	}

	if err := h.financeService.UpsertBudget(budget); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, budget)
}

func (h *FinanceHandler) GetDashboard(c *gin.Context) {
	userID := getUserID(c)
	summary, err := h.financeService.GetDashboardSummary(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, summary)
}

func getUserID(c *gin.Context) uuid.UUID {
	userIDstr, _ := c.Get("user_id")
	id, _ := uuid.Parse(userIDstr.(string))
	return id
}
