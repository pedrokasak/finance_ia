package user

import "github.com/google/uuid"

type Repository interface {
	Create(user *User) error
	FindAll() ([]*User, error)
	FindByID(id uuid.UUID) (*User, error)
	FindByEmail(email string) (*User, error)
	Update(user *User) error
	Delete(user *User) error
}
