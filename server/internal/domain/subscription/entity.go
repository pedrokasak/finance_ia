package subscription

import (
	"time"

	"github.com/google/uuid"
)

// SubscriptionStatus tracks the Stripe subscription lifecycle
type SubscriptionStatus string

const (
	StatusActive   SubscriptionStatus = "active"
	StatusCanceled SubscriptionStatus = "canceled"
	StatusPastDue  SubscriptionStatus = "past_due"
	StatusTrialing SubscriptionStatus = "trialing"
)

// Subscription tracks the user's payment subscription
type Subscription struct {
	ID                 uuid.UUID          `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID             uuid.UUID          `json:"user_id" gorm:"type:uuid;not null;uniqueIndex"`
	Plan               string             `json:"plan" gorm:"not null;default:'free'"`
	Status             SubscriptionStatus `json:"status" gorm:"not null;default:'active'"`
	ExternalID         string             `json:"external_id" gorm:"default:''"` // Stripe subscription ID
	ExternalCustomerID string             `json:"-" gorm:"default:''"`           // Stripe customer ID
	PriceID            string             `json:"-" gorm:"default:''"`           // Stripe price ID
	CurrentPeriodStart *time.Time         `json:"current_period_start"`
	CurrentPeriodEnd   *time.Time         `json:"current_period_end"`
	CanceledAt         *time.Time         `json:"canceled_at,omitempty"`
	CreatedAt          time.Time          `json:"created_at"`
	UpdatedAt          time.Time          `json:"updated_at"`
}

// CheckoutSession holds the result of creating a Stripe checkout session
type CheckoutSession struct {
	SessionID string `json:"session_id"`
	URL       string `json:"url"`
}

// PortalSession holds the result of creating a Stripe billing portal session
type PortalSession struct {
	URL string `json:"url"`
}

// PlanFeatures defines what each plan unlocks
type PlanFeatures struct {
	Plan              string   `json:"plan"`
	AIInsightsPerWeek int      `json:"ai_insights_per_week"` // -1 = unlimited
	FullAnalysis      bool     `json:"full_analysis"`
	Projections       bool     `json:"projections"`
	SmartAlerts       bool     `json:"smart_alerts"`
	ExportData        bool     `json:"export_data"`
	MaxTransactions   int      `json:"max_transactions"` // -1 = unlimited
	Features          []string `json:"features"`
}

// GetPlanFeatures returns the feature set for a given plan
func GetPlanFeatures(plan string) PlanFeatures {
	switch plan {
	case "premium":
		return PlanFeatures{
			Plan:              "premium",
			AIInsightsPerWeek: -1,
			FullAnalysis:      true,
			Projections:       true,
			SmartAlerts:       true,
			ExportData:        true,
			MaxTransactions:   -1,
			Features: []string{
				"IA financeira completa",
				"Análise mensal detalhada",
				"Diagnóstico de saúde financeira",
				"Projeções financeiras",
				"Alertas inteligentes",
				"Exportação de dados",
				"Transações ilimitadas",
			},
		}
	case "pro":
		return PlanFeatures{
			Plan:              "pro",
			AIInsightsPerWeek: -1,
			FullAnalysis:      true,
			Projections:       true,
			SmartAlerts:       true,
			ExportData:        true,
			MaxTransactions:   -1,
			Features: []string{
				"Tudo do Premium",
				"Copiloto financeiro em tempo real",
				"Score de saúde gamificado",
				"Planejamento de longo prazo",
				"Simulação de independência financeira",
				"Análise de comportamento financeiro",
				"Suporte prioritário 24/7",
			},
		}
	default: // free
		return PlanFeatures{
			Plan:              "free",
			AIInsightsPerWeek: 1,
			FullAnalysis:      false,
			Projections:       false,
			SmartAlerts:       false,
			ExportData:        false,
			MaxTransactions:   100,
			Features: []string{
				"Cadastro de receitas e despesas",
				"Dashboard básico",
				"1 insight de IA por semana",
				"Alertas simples de risco",
				"Categorização manual",
			},
		}
	}
}
