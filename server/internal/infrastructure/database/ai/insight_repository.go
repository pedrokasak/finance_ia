package ai

import (
	"finance-ia/internal/domain/ai"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InsightRepository struct {
	db *gorm.DB
}

func NewInsightRepository(db *gorm.DB) *InsightRepository {
	return &InsightRepository{db: db}
}

func (r *InsightRepository) Save(insight *ai.AIInsight) error {
	return r.db.Create(insight).Error
}

func (r *InsightRepository) FindLatestByUser(userID uuid.UUID, period string) (*ai.AIInsight, error) {
	var insight ai.AIInsight
	if err := r.db.
		Where("user_id = ? AND period = ?", userID, period).
		Order("generated_at DESC").
		First(&insight).Error; err != nil {
		return nil, err
	}
	return &insight, nil
}

func (r *InsightRepository) CountByUserAndPeriod(userID uuid.UUID, start, end time.Time) (int64, error) {
	var count int64
	if err := r.db.Model(&ai.AIInsight{}).
		Where("user_id = ? AND generated_at BETWEEN ? AND ?", userID, start, end).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
