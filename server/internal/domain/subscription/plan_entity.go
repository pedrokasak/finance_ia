package subscription

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Plan represents a subscription plan available for purchase.
// Plans are seeded into the database and linked to Stripe Price IDs.
type Plan struct {
	ID                   uuid.UUID     `json:"id" gorm:"type:uuid;primaryKey"`
	Slug                 string        `json:"slug" gorm:"unique;not null"` // "free", "pro", "premium"
	Name                 string        `json:"name" gorm:"not null"`        // Display name
	Description          string        `json:"description"`
	PriceMonthly         float64       `json:"price_monthly" gorm:"default:0"`
	PriceYearly          float64       `json:"price_yearly" gorm:"default:0"`
	StripeProductID      string        `json:"stripe_product_id" gorm:"default:''"`
	StripePriceIDMonthly string        `json:"stripe_price_id_monthly" gorm:"default:''"`
	StripePriceIDYearly  string        `json:"stripe_price_id_yearly" gorm:"default:''"`
	Features             []PlanFeature `json:"features" gorm:"foreignKey:PlanID;constraint:OnDelete:CASCADE"`
	MaxTransactions      int           `json:"max_transactions" gorm:"default:100"`
	AIInsights           bool          `json:"ai_insights" gorm:"default:false"`
	ExportData           bool          `json:"export_data" gorm:"default:false"`
	IsActive             bool          `json:"is_active" gorm:"default:true"`
}

func (p *Plan) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// PlanFeature represents a specific feature of a subscription plan.
type PlanFeature struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	PlanID      uuid.UUID `json:"plan_id" gorm:"type:uuid;index"`
	Description string    `json:"description" gorm:"not null"`
}

func (pf *PlanFeature) BeforeCreate(tx *gorm.DB) error {
	if pf.ID == uuid.Nil {
		pf.ID = uuid.New()
	}
	return nil
}

// PlanRepository defines the data access contract for plans
type PlanRepository interface {
	FindAll() ([]*Plan, error)
	FindBySlug(slug string) (*Plan, error)
	FindByID(id uuid.UUID) (*Plan, error)
	FindByStripePriceID(priceID string) (*Plan, error)
	Upsert(plan *Plan) error
	Delete(id uuid.UUID) error

	UpsertFeature(feature *PlanFeature) error
	DeleteFeature(id uuid.UUID) error
}
