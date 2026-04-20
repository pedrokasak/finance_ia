package handlers

import (
	"errors"
	"finance-ia/internal/application/usecase/user"
	"finance-ia/internal/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserHandler struct {
	usecase user.IUserUseCase
}

func NewUserHandler(uc user.IUserUseCase) *UserHandler {
	return &UserHandler{usecase: uc}
}

func (h *UserHandler) RegisterRoutes(public, protected gin.IRouter) {
	public.POST("/user/register", h.Register)
	protected.GET("/user/:id", h.GetUserByID)
	protected.PUT("/user/update/:id", h.UpdateUser)
	protected.DELETE("/user/delete/:id", h.DeleteUser)
}

func (h *UserHandler) Register(c *gin.Context) {
	var req struct {
		FirstName string `json:"first_name" binding:"required,min=1,max=80"`
		LastName  string `json:"last_name" binding:"required,min=1,max=80"`
		Email     string `json:"email" binding:"required,email,max=254"`
		Password  string `json:"password" binding:"required,min=8,max=72"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos"})
		return
	}

	user, err := h.usecase.Register(strings.TrimSpace(req.FirstName), strings.TrimSpace(req.LastName), strings.TrimSpace(strings.ToLower(req.Email)), req.Password)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			c.JSON(http.StatusConflict, gin.H{"error": "email already in use"})
			return
		}
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
	currentUserID := utils.GetUserID(c)
	if currentUserID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	if currentUserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "acesso negado"})
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
	currentUserID := utils.GetUserID(c)
	if currentUserID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	if currentUserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "acesso negado"})
		return
	}

	var req struct {
		ID                   uuid.UUID `json:"id"`
		Email                string    `json:"email" binding:"omitempty,email,max=254"`
		FirstName            string    `json:"first_name" binding:"omitempty,max=80"`
		LastName             string    `json:"last_name" binding:"omitempty,max=80"`
		AvatarURL            *string   `json:"avatar_url"`
		FinancialMethodID    *string   `json:"financial_method_id"`
		NotificationsEnabled *bool     `json:"notifications_enabled"`
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
		userObj.Email = strings.TrimSpace(strings.ToLower(req.Email))
	}

	if req.FirstName != "" {
		userObj.FirstName = strings.TrimSpace(req.FirstName)
	}
	if req.LastName != "" {
		userObj.LastName = strings.TrimSpace(req.LastName)
	}
	if req.AvatarURL != nil {
		avatarURL := strings.TrimSpace(*req.AvatarURL)
		if err := utils.ValidateAvatarDataURL(avatarURL); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userObj.AvatarURL = avatarURL
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

	if err := h.usecase.Update(userObj); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao atualizar perfil"})
		return
	}
	c.JSON(http.StatusOK, userObj)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	currentUserID := utils.GetUserID(c)
	if currentUserID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	if currentUserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "acesso negado"})
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
