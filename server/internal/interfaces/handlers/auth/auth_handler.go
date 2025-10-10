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
}

func (h *AuthHandler) Login(c *gin.Context) {
    var req struct {
        Email    string `json:"email" binding:"required,email"`
        Password string `json:"password" binding:"required"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid data", "error": err.Error()})
        return
    }

    token, err := h.usecase.Login(req.Email, req.Password)
    if err != nil {
        // mapear erros do usecase para status apropriado
        c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials", "error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "email":   req.Email,
        "token":   token,
        "success": true,
    })
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
    var req struct {
        Email string `json:"email" binding:"required,email"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid data", "error": err.Error()})
        return
    }

    if err := h.usecase.ForgotPassword(req.Email); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Error processing request"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "If the email exists, a reset link has been sent to it."})
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
    var req struct {
        Token       string `json:"token" binding:"required"`
        NewPassword string `json:"new_password" binding:"required,min=6"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid data", "error": err.Error()})
        return
    }

    if err := h.usecase.ResetPassword(req.Token, req.NewPassword); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid or expired token", "error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}
