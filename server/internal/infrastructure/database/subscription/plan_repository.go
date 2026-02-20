package subscription

import (
	"finance-ia/internal/domain/subscription"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PlanRepository struct {
	db *gorm.DB
}

func NewPlanRepository(db *gorm.DB) *PlanRepository {
	return &PlanRepository{db: db}
}

func (r *PlanRepository) FindAll() ([]*subscription.Plan, error) {
	var plans []*subscription.Plan
	if err := r.db.Preload("Features").Where("is_active = ?", true).Order("price_monthly ASC").Find(&plans).Error; err != nil {
		return nil, err
	}
	return plans, nil
}

func (r *PlanRepository) FindBySlug(slug string) (*subscription.Plan, error) {
	var plan subscription.Plan
	if err := r.db.Preload("Features").Where("slug = ?", slug).First(&plan).Error; err != nil {
		return nil, err
	}
	return &plan, nil
}

func (r *PlanRepository) FindByID(id uuid.UUID) (*subscription.Plan, error) {
	var plan subscription.Plan
	if err := r.db.Preload("Features").Where("id = ?", id).First(&plan).Error; err != nil {
		return nil, err
	}
	return &plan, nil
}

func (r *PlanRepository) Upsert(plan *subscription.Plan) error {
	return r.db.Save(plan).Error
}

func (r *PlanRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&subscription.Plan{}, id).Error
}

func (r *PlanRepository) UpsertFeature(feature *subscription.PlanFeature) error {
	return r.db.Save(feature).Error
}

func (r *PlanRepository) DeleteFeature(id uuid.UUID) error {
	return r.db.Delete(&subscription.PlanFeature{}, id).Error
}

var _ subscription.PlanRepository = (*PlanRepository)(nil)
