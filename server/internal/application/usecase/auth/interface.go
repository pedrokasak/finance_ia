package auth

import (
	"finance-ia/internal/domain/auth"
)

type IAuthUseCase interface {
	Login(email, password, code string) (string, error)
	Setup2FA(email string) (string, string, error)
	Verify2FA(email, code string) error
	Disable2FA(email string) error
	GetByEmail(email string) (*auth.Authentication, error)
	ForgotPassword(email string) error
	ResetPassword(token, newPassword string) error
	Update(user *auth.Authentication) error
	Logout(token string) error
}

var _ IAuthUseCase = &UseCase{}
