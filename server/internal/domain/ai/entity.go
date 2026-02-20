package ai

import (
	"time"

	"github.com/google/uuid"
)

// InsightType categorizes the AI insight
type InsightType string

const (
	InsightTypeWarning     InsightType = "warning"
	InsightTypeTip         InsightType = "tip"
	InsightTypeAchievement InsightType = "achievement"
	InsightTypeProjection  InsightType = "projection"
)

// AIInsight represents a generated AI financial insight for a user
type AIInsight struct {
	ID          uuid.UUID   `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID      uuid.UUID   `json:"user_id" gorm:"type:uuid;not null;index"`
	Type        InsightType `json:"type" gorm:"not null"`
	Title       string      `json:"title" gorm:"not null"`
	Content     string      `json:"content" gorm:"not null"`
	Plan        string      `json:"plan" gorm:"not null"` // which plan generated this
	GeneratedAt time.Time   `json:"generated_at"`
	ExpiresAt   time.Time   `json:"expires_at"`
	Period      string      `json:"period" gorm:"not null"` // "2024-01"
}

// FinancialContext is passed to the AI provider for analysis
type FinancialContext struct {
	UserName      string
	Plan          string
	TotalIncome   float64
	TotalExpenses float64
	Balance       float64
	SavingsRate   float64
	HealthScore   int
	HealthLevel   string
	TopCategories []CategorySpend
	MonthlyTrends []MonthTrend
	Period        string
}

type CategorySpend struct {
	Name       string
	Amount     float64
	Percentage float64
}

type MonthTrend struct {
	Month    string
	Income   float64
	Expenses float64
}
