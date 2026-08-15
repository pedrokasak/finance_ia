package user

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Plan represents the user's subscription plan
type Plan string

const (
	PlanFree    Plan = "free"
	PlanPremium Plan = "premium"
	PlanPro     Plan = "pro"
)

// FinancialMethod represents the budgeting method chosen by the user
type FinancialMethod string

const (
	MethodFiftyThirtyTwenty FinancialMethod = "50-30-20"
	MethodEnvelopes         FinancialMethod = "envelopes"
	MethodZeroBased         FinancialMethod = "zero-based"
	MethodSeventyTwentyTen  FinancialMethod = "70-20-10"
)

type User struct {
	ID                   uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey"`
	FirstName            string     `json:"first_name" gorm:"not null"`
	LastName             string     `json:"last_name" gorm:"not null"`
	Email                string     `json:"email" gorm:"unique;not null"`
	Password             string     `json:"-" gorm:"not null"`
	AvatarURL            string     `json:"avatar_url" gorm:"default:''"`
	Plan                 Plan       `json:"plan" gorm:"default:'free';not null"`
	FinancialMethodID    *uuid.UUID `json:"financial_method_id" gorm:"type:uuid;index"`
	MonthlyIncome        float64    `json:"monthly_income" gorm:"default:0"`
	OnboardingCompleted  bool       `json:"onboarding_completed" gorm:"default:false"`
	NotificationsEnabled bool       `json:"notifications_enabled" gorm:"default:true"`
	TwoFAEnabled         bool       `json:"two_fa_enabled" gorm:"default:false"`
	TwoFASecret          string     `json:"-" gorm:"default:''"`
	StripeCustomerID     string     `json:"-" gorm:"default:''"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
