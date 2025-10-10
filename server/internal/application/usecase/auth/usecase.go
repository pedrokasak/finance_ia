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

func (uc *UseCase) Login(email, password string) (string, error) {
	return uc.service.Login(email, password)
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
