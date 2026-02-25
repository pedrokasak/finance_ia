package handlers

import (
	"finance-ia/internal/application/usecase/user"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	usecase user.IUserUseCase
}

func NewUserHandler(uc user.IUserUseCase) *UserHandler {
	return &UserHandler{usecase: uc}
}

func (h *UserHandler) RegisterRoutes(public, protected gin.IRouter) {
	public.POST("/user/register", h.Register)
	protected.GET("/users/", h.GetAllUsers)
	protected.GET("/user/:id", h.GetUserByID)
	protected.PUT("/user/update/:id", h.UpdateUser)
	protected.DELETE("/user/delete/:id", h.DeleteUser)
}

func (h *UserHandler) Register(c *gin.Context) {
	var req struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos"})
		return
	}

	_, err := h.usecase.GetByEmail(req.Email)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "email already in use"})
		return
	}

	user, err := h.usecase.Register(req.FirstName, req.LastName, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.usecase.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao buscar usuários"})
		return
	}
	if len(users) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "nenhum usuário encontrado"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	user, err := h.usecase.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var req struct {
		ID                   uuid.UUID `json:"id"`
		Email                string    `json:"email"`
		FirstName            string    `json:"first_name"`
		LastName             string    `json:"last_name"`
		AvatarURL            *string   `json:"avatar_url"`
		FinancialMethodID    *string   `json:"financial_method_id"`
		NotificationsEnabled *bool     `json:"notifications_enabled"`
		MonthlyIncome        *float64  `json:"monthly_income"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos"})
		return
	}

	userObj, err := h.usecase.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "usuário não encontrado"})
		return
	}

	if req.Email != "" {
		userObj.Email = req.Email
	}

	if req.FirstName != "" {
		userObj.FirstName = req.FirstName
	}
	if req.LastName != "" {
		userObj.LastName = req.LastName
	}
	if req.AvatarURL != nil {
		userObj.AvatarURL = *req.AvatarURL
	}
	if req.NotificationsEnabled != nil {
		userObj.NotificationsEnabled = *req.NotificationsEnabled
	}
	if req.FinancialMethodID != nil {
		if *req.FinancialMethodID == "" {
			userObj.FinancialMethodID = nil
		} else {
			methodUUID, err := uuid.Parse(*req.FinancialMethodID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid financial_method_id format"})
				return
			}
			userObj.FinancialMethodID = &methodUUID
		}
	}
	if req.MonthlyIncome != nil {
		userObj.MonthlyIncome = *req.MonthlyIncome
	}

	if err := h.usecase.Update(userObj); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao atualizar perfil"})
		return
	}
	c.JSON(http.StatusOK, userObj)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	userObj, err := h.usecase.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "usuário não encontrado"})
		return
	}

	if err := h.usecase.Delete(userObj); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao deletar usuário"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully", "user": userObj})
}
