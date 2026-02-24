package goal

import (
	"errors"

	"github.com/google/uuid"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(goal *Goal) error {
	if goal.TargetAmount <= 0 {
		return errors.New("target_amount must be greater than zero")
	}
	if goal.Name == "" {
		return errors.New("name is required")
	}
	return s.repo.Create(goal)
}

func (s *Service) GetByID(id uuid.UUID) (*Goal, error) {
	return s.repo.FindByID(id)
}

func (s *Service) GetByUser(userID uuid.UUID) ([]*Goal, error) {
	return s.repo.FindByUserID(userID)
}

func (s *Service) Update(goal *Goal) error {
	if goal.CurrentAmount < 0 {
		return errors.New("current_amount cannot be negative")
	}
	return s.repo.Update(goal)
}

func (s *Service) Delete(id uuid.UUID) error {
	goal, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(goal)
}
