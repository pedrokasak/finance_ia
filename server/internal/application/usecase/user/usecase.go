package user

import "finance-ia/internal/domain/user"

type UseCase struct {
	service *user.Service
}

func NewUseCase(service *user.Service) *UseCase {
	return &UseCase{service: service}
}

func (uc *UseCase) Register(firstName, LastName, email, password string) (*user.User, error) {
	return uc.service.Register(firstName, LastName, email, password)
}

func (uc *UseCase) Login(email, password string) (string, error) {
	return uc.service.Login(email, password)
}

func (uc *UseCase) GetProfile(userID uint) (*user.User, error) {
	return uc.service.GetByID(userID)
}

func (uc *UseCase) GetByEmail(email string) (*user.User, error) {
	return uc.service.GetByEmail(email)
}

func (uc *UseCase) GetByID(id uint) (*user.User, error) {
	return uc.service.GetByID(id)
}

func (uc *UseCase) Update(user *user.User) error {
	return uc.service.Update(user)
}

func (uc *UseCase) Delete(user *user.User) error {
	return uc.service.Delete(user)
}