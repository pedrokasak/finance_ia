package auth

import (
	"finance-ia/internal/domain/auth"
)


type IAuthUseCase interface {
	Login(email, password string) (string, error)
	GetByEmail(email string) (*auth.Authentication, error)
	ForgotPassword(email string) error
	ResetPassword(token, newPassword string) error
	Update(user *auth.Authentication) error
	Logout(token string) error
}

var _ IAuthUseCase = &UseCase{}