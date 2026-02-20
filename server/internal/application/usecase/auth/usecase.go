package auth

import (
	"finance-ia/internal/domain/auth"
)

type UseCase struct {
	service *auth.Service
}

func NewUseCase(service *auth.Service) *UseCase {
	return &UseCase{service: service}
}

func (uc *UseCase) Login(email, password, code string) (string, error) {
	return uc.service.Login(email, password, code)
}

func (uc *UseCase) Setup2FA(email string) (string, string, error) {
	return uc.service.Setup2FA(email)
}

func (uc *UseCase) Verify2FA(email, code string) error {
	return uc.service.Verify2FA(email, code)
}

func (uc *UseCase) Disable2FA(email string) error {
	return uc.service.Disable2FA(email)
}

func (uc *UseCase) ForgotPassword(email string) error {
	return uc.service.ForgotPassword(email)
}

func (uc *UseCase) ResetPassword(token, newPassword string) error {
	return uc.service.ResetPassword(token, newPassword)
}

func (uc *UseCase) GetByEmail(email string) (*auth.Authentication, error) {
	return uc.service.GetByEmail(email)
}

func (uc *UseCase) Update(user *auth.Authentication) error {
	return uc.service.Update(user)
}

func (uc *UseCase) Logout(token string) error {
	return uc.service.Logout(token)
}
