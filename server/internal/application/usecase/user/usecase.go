package user

import (
	"finance-ia/internal/domain/user"

	"github.com/google/uuid"
)

type UseCase struct {
	service *user.Service
}

func NewUseCase(service *user.Service) *UseCase {
	return &UseCase{service: service}
}

func (uc *UseCase) Register(firstName, LastName, email, password string) (*user.User, error) {
	return uc.service.Register(firstName, LastName, email, password)
}

func (uc *UseCase) GetAll() ([]*user.User, error) {
	return uc.service.GetAll()
}

func (uc *UseCase) GetByID(userID uuid.UUID) (*user.User, error) {
	return uc.service.GetByID(userID)
}

func (uc *UseCase) GetByEmail(email string) (*user.User, error) {
	return uc.service.GetByEmail(email)
}

func (uc *UseCase) Update(user *user.User) error {
	return uc.service.Update(user)
}

func (uc *UseCase) Delete(user *user.User) error {
	return uc.service.Delete(user)
}