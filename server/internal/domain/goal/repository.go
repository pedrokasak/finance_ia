package goal

import "github.com/google/uuid"

type Repository interface {
	Create(goal *Goal) error
	FindByID(id uuid.UUID) (*Goal, error)
	FindByUserID(userID uuid.UUID) ([]*Goal, error)
	Update(goal *Goal) error
	Delete(goal *Goal) error
}
