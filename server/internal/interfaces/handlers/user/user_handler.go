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
    userIDStr := c.MustGet("id").(string)
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
	var req struct {
		ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}

	userID := c.MustGet("id").(string)
	user, err := h.usecase.GetByEmail(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "usuário não encontrado"})
		return
	}

	user.ID = req.ID

	if err := h.usecase.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao atualizar perfil"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	userIDStr := c.MustGet("id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	user, err := h.usecase.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "usuário não encontrado"})
		return
	}

	if err := h.usecase.Delete(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao deletar usuário"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully", "user": user})
}

