package finance

import (
	"time"

	"github.com/google/uuid"
)

// TransactionType represents whether a transaction is income or expense
type TransactionType string

const (
	TransactionTypeIncome  TransactionType = "income"
	TransactionTypeExpense TransactionType = "expense"
)

// Transaction represents a financial transaction
type Transaction struct {
	ID             uuid.UUID       `json:"id" gorm:"type:uuid;primary_key"`
	UserID         uuid.UUID       `json:"user_id" gorm:"type:uuid;not null;index"`
	CategoryID     *uuid.UUID      `json:"category_id" gorm:"type:uuid"`
	Type           TransactionType `json:"type" gorm:"not null"`
	Amount         float64         `json:"amount" gorm:"not null"`
	Description    string          `json:"description" gorm:"default:''"`
	Date           time.Time       `json:"date" gorm:"not null"`
	IsRecurring    bool            `json:"is_recurring" gorm:"default:false"`
	IdempotencyKey string          `json:"-" gorm:"uniqueIndex;default:''"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`

	// Associations
	Category *Category `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
}

// CategoryType groups categories as income or expense
type CategoryType string

const (
	CategoryTypeIncome  CategoryType = "income"
	CategoryTypeExpense CategoryType = "expense"
)

// Category represents an expense or income category
type Category struct {
	ID        uuid.UUID    `json:"id" gorm:"type:uuid;primary_key"`
	UserID    *uuid.UUID   `json:"user_id" gorm:"type:uuid"` // nil = default category
	Name      string       `json:"name" gorm:"not null"`
	Type      CategoryType `json:"type" gorm:"not null"`
	Color     string       `json:"color" gorm:"default:'#6366f1'"`
	Icon      string       `json:"icon" gorm:"default:'tag'"`
	IsDefault bool         `json:"is_default" gorm:"default:false"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

// Budget stores the user's monthly budget allocation by financial method
type Budget struct {
	ID             uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	UserID         uuid.UUID `json:"user_id" gorm:"type:uuid;not null;uniqueIndex:idx_budget_user_period"`
	Period         string    `json:"period" gorm:"not null;uniqueIndex:idx_budget_user_period"` // "2024-01"
	TotalIncome    float64   `json:"total_income" gorm:"not null"`
	NeedsPercent   float64   `json:"needs_percent" gorm:"default:50"`   // Necessidades
	WantsPercent   float64   `json:"wants_percent" gorm:"default:30"`   // Desejos
	SavingsPercent float64   `json:"savings_percent" gorm:"default:20"` // Poupança/Investimento
	NeedsAmount    float64   `json:"needs_amount"`
	WantsAmount    float64   `json:"wants_amount"`
	SavingsAmount  float64   `json:"savings_amount"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// DashboardSummary aggregates financial data for the dashboard
type DashboardSummary struct {
	TotalIncome       float64           `json:"total_income"`
	TotalExpenses     float64           `json:"total_expenses"`
	Balance           float64           `json:"balance"`
	SavingsRate       float64           `json:"savings_rate"` // Percentage
	HealthScore       int               `json:"health_score"` // 0-1000
	HealthLevel       string            `json:"health_level"` // Bronze, Prata, Ouro, Platina, Diamante
	Budget            *Budget           `json:"budget,omitempty"`
	CategoryBreakdown []CategorySummary `json:"category_breakdown"`
	MonthlyTrend      []MonthlyTrend    `json:"monthly_trend"`
	DaysUntilNegative *int              `json:"days_until_negative,omitempty"` // risk alert
}

// CategorySummary holds aggregated data per category
type CategorySummary struct {
	CategoryID   uuid.UUID `json:"category_id"`
	CategoryName string    `json:"category_name"`
	Color        string    `json:"color"`
	Total        float64   `json:"total"`
	Percentage   float64   `json:"percentage"`
}

// MonthlyTrend holds monthly income vs expense data for charts
type MonthlyTrend struct {
	Month    string  `json:"month"`
	Income   float64 `json:"income"`
	Expenses float64 `json:"expenses"`
}

// BehavioralAnalysis holds advanced behavioral patterns detected in transactions
type BehavioralAnalysis struct {
	// EmotionalSpendingScore: 0-100 — % of expenses on weekends or late hours (impulse indicator)
	EmotionalSpendingScore float64 `json:"emotional_spending_score"`

	// SalaryEffectDays: average extra spending in the N days after income received (vs normal days)
	SalaryEffectDays   int     `json:"salary_effect_days"`   // how many days the "effect" lasts
	SalaryEffectAmount float64 `json:"salary_effect_amount"` // average extra spend per day after salary

	// WeekendVsWeekdayRatio: ratio of avg weekend spend to avg weekday spend (>1.5 = risk)
	WeekendVsWeekdayRatio float64 `json:"weekend_vs_weekday_ratio"`

	// MoneyLeakCategories: categories with consistent month-over-month growth > 10%
	MoneyLeakCategories []MoneyLeakCategory `json:"money_leak_categories"`

	// ImpulsivePurchaseCount: transactions < R$50 and non-recurring in the last 30 days
	ImpulsivePurchaseCount int     `json:"impulsive_purchase_count"`
	ImpulsivePurchaseTotal float64 `json:"impulsive_purchase_total"`

	// TopWeekdayByExpense: day of week with highest avg spending (0=Sun, 6=Sat)
	TopWeekdayByExpense int    `json:"top_weekday_by_expense"`
	TopWeekdayName      string `json:"top_weekday_name"`

	// RiskLevel: "low", "medium", "high" based on combined score
	RiskLevel   string `json:"risk_level"`
	RiskMessage string `json:"risk_message"`
}

// MoneyLeakCategory is a category showing consistent spending growth
type MoneyLeakCategory struct {
	CategoryName  string  `json:"category_name"`
	CategoryColor string  `json:"category_color"`
	GrowthPercent float64 `json:"growth_percent"` // month-over-month avg growth
	CurrentMonth  float64 `json:"current_month"`
	PreviousMonth float64 `json:"previous_month"`
}
