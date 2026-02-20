package finance

import (
	"finance-ia/internal/domain/finance"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BudgetRepository struct {
	db *gorm.DB
}

func NewBudgetRepository(db *gorm.DB) *BudgetRepository {
	return &BudgetRepository{db: db}
}

func (r *BudgetRepository) Upsert(budget *finance.Budget) error {
	return r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "user_id"}, {Name: "period"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"total_income", "needs_percent", "wants_percent", "savings_percent",
			"needs_amount", "wants_amount", "savings_amount", "updated_at",
		}),
	}).Create(budget).Error
}

func (r *BudgetRepository) FindByUserAndPeriod(userID uuid.UUID, period string) (*finance.Budget, error) {
	var budget finance.Budget
	if err := r.db.Where("user_id = ? AND period = ?", userID, period).First(&budget).Error; err != nil {
		return nil, err
	}
	return &budget, nil
}
