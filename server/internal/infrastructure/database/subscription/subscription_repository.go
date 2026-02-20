package subscription

import (
	"finance-ia/internal/domain/subscription"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SubscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) Upsert(sub *subscription.Subscription) error {
	return r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"plan", "status", "external_id", "external_customer_id",
			"price_id", "current_period_start", "current_period_end",
			"canceled_at", "updated_at",
		}),
	}).Create(sub).Error
}

func (r *SubscriptionRepository) FindByUserID(userID uuid.UUID) (*subscription.Subscription, error) {
	var sub subscription.Subscription
	if err := r.db.Where("user_id = ?", userID).First(&sub).Error; err != nil {
		return nil, err
	}
	return &sub, nil
}

func (r *SubscriptionRepository) FindByExternalID(externalID string) (*subscription.Subscription, error) {
	var sub subscription.Subscription
	if err := r.db.Where("external_id = ?", externalID).First(&sub).Error; err != nil {
		return nil, err
	}
	return &sub, nil
}
