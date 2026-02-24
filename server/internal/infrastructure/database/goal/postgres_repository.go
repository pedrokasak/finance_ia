package goal

import (
	"finance-ia/internal/domain/goal"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostgresGoalRepository struct {
	db *gorm.DB
}

func NewPostgresGoalRepository(db *gorm.DB) *PostgresGoalRepository {
	return &PostgresGoalRepository{db: db}
}

func (r *PostgresGoalRepository) Create(goal *goal.Goal) error {
	return r.db.Create(goal).Error
}

func (r *PostgresGoalRepository) FindByID(id uuid.UUID) (*goal.Goal, error) {
	var g goal.Goal
	if err := r.db.Where("id = ?", id).First(&g).Error; err != nil {
		return nil, err
	}
	return &g, nil
}

func (r *PostgresGoalRepository) FindByUserID(userID uuid.UUID) ([]*goal.Goal, error) {
	var goals []*goal.Goal
	if err := r.db.Where("user_id = ?", userID).Find(&goals).Error; err != nil {
		return nil, err
	}
	return goals, nil
}

func (r *PostgresGoalRepository) Update(goal *goal.Goal) error {
	return r.db.Save(goal).Error
}

func (r *PostgresGoalRepository) Delete(goal *goal.Goal) error {
	return r.db.Delete(goal).Error
}
