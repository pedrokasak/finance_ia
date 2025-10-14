package auth

import (
	"finance-ia/internal/application/usecase/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)


type AuthHandler struct {
    usecase auth.IAuthUseCase
}

func NewAuthHandler(uc auth.IAuthUseCase) *AuthHandler {
    return &AuthHandler{usecase: uc}
}

func (h *AuthHandler) RegisterRoutes(public, protected gin.IRouter) {
    public.POST("/auth/login", h.Login)
    public.POST("/auth/forgot-password", h.ForgotPassword)
    public.POST("/auth/reset-password", h.ResetPassword)
    protected.POST("/auth/logout", h.Logout)
}

func (h *AuthHandler) Login(c *gin.Context) {
    var req struct {
        Email    string `json:"email" binding:"required,email"`
        Password string `json:"password" binding:"required"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "Invalid data"})
        return
    }

    token, err := h.usecase.Login(req.Email, req.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error(), "message": "Not authorized"})
        return
    }

    if token == "" {
        c.JSON(http.StatusForbidden, gin.H{"message": "Invalid credentials"})
        return
    }
    

    c.JSON(http.StatusOK, gin.H{
        "email":   req.Email,
        "token":   token,
        "success": true,
    })
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   err.Error(),
			Message: "Dados inválidos",
		})
		return
	}

	if err := h.usecase.ForgotPassword(req.Email); err != nil {
		// Log o erro mas retorne mensagem genérica
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: "Erro ao processar solicitação",
		})
		return
	}

	c.JSON(http.StatusOK, MessageResponse{
		Message: "Se o email existir, você receberá um link de recuperação.",
	})
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   err.Error(),
			Message: "Dados inválidos",
		})
		return
	}

	if err := h.usecase.ResetPassword(req.Token, req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
			Error:   "Token inválido ou expirado",
		})
		return
	}

	c.JSON(http.StatusOK, MessageResponse{
		Message: "Senha redefinida com sucesso!",
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
    authHeader := c.GetHeader("Authorization")
    if authHeader == "" {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Authorization header required"})
        return
    }

    token := authHeader[len("Bearer "):]
    if err := h.usecase.Logout(token); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Error logging out", "error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}


type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}