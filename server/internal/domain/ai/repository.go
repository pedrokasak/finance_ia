package ai

import (
	"time"

	"github.com/google/uuid"
)

// InsightRepository defines persistence for AI insights
type InsightRepository interface {
	Save(insight *AIInsight) error
	FindLatestByUser(userID uuid.UUID, period string) (*AIInsight, error)
	CountByUserAndPeriod(userID uuid.UUID, start, end time.Time) (int64, error)
}

// AIProvider is the port/interface for AI providers (Gemini, OpenAI, etc.)
// This is the adapter pattern — swap providers without touching business logic
type AIProvider interface {
	// GenerateInsight generates a short financial insight (used for free plan)
	GenerateInsight(ctx FinancialContext) (*AIInsight, error)
	// GenerateFullAnalysis generates a comprehensive monthly analysis (premium+)
	GenerateFullAnalysis(ctx FinancialContext) ([]*AIInsight, error)
}
