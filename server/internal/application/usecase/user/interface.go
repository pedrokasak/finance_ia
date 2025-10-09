package user

import (
	"finance-ia/internal/domain/user"

	"github.com/google/uuid"
)


type IUserUseCase interface {
	Register(firstName, lastName, email, password string) (*user.User, error)
	Login(email, password string) (string, error)
	GetAll() ([]*user.User, error)
	GetByID(userID uuid.UUID) (*user.User, error)
	GetByEmail(email string) (*user.User, error)
	Update(user *user.User) error
	Delete(user *user.User) error
	ForgotPassword(email string) error
	ResetPassword(token, newPassword string) error
}

var _ IUserUseCase = &UseCase{}