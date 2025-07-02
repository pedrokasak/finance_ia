package handlers

import (
	"finance-ia/internal/application/user"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	usecase *user.UseCase
}

func NewUserHandler(uc *user.UseCase) *UserHandler {
	return &UserHandler{usecase: uc}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req struct {
		FirstName string `json:"first_name"`
		LastName	string `json:"last_name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos"})
		return
	}
	user, err := h.usecase.Register(req.FirstName, req.LastName, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}
 func (h *UserHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos"})
		return
	}
	// Validate email and password
	token, err := h.usecase.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"email": req.Email,
		"password":  req.Password,
		"token": token,
		"success": true,
	},
	)
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	user, err := h.usecase.GetByEmail(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "usuário não encontrado"})
		return
	}
	c.JSON(http.StatusOK, user)
}
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	var req struct {
		ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos"})
		return
	}

	userID := c.MustGet("user_id").(string)
	user, err := h.usecase.GetByEmail(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "usuário não encontrado"})
		return
	}

	user.ID = req.ID

	if err := h.usecase.UpdateProfile(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao atualizar perfil"})
		return
	}

	c.JSON(http.StatusOK, user)
}
