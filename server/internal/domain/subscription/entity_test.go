package subscription_test

import (
	"finance-ia/internal/domain/subscription"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPlanFeatures(t *testing.T) {
	t.Run("free plan", func(t *testing.T) {
		features := subscription.GetPlanFeatures("free")
		assert.Equal(t, "free", features.Plan)
		assert.Equal(t, 1, features.AIInsightsPerWeek)
		assert.False(t, features.FullAnalysis)
		assert.False(t, features.Projections)
		assert.False(t, features.SmartAlerts)
		assert.False(t, features.ExportData)
		assert.Equal(t, 100, features.MaxTransactions)
		assert.Len(t, features.Features, 5)
	})

	t.Run("premium plan", func(t *testing.T) {
		features := subscription.GetPlanFeatures("premium")
		assert.Equal(t, "premium", features.Plan)
		assert.Equal(t, -1, features.AIInsightsPerWeek)
		assert.True(t, features.FullAnalysis)
		assert.True(t, features.Projections)
		assert.True(t, features.SmartAlerts)
		assert.True(t, features.ExportData)
		assert.Equal(t, -1, features.MaxTransactions)
		assert.Len(t, features.Features, 7)
	})

	t.Run("pro plan", func(t *testing.T) {
		features := subscription.GetPlanFeatures("pro")
		assert.Equal(t, "pro", features.Plan)
		assert.Equal(t, -1, features.AIInsightsPerWeek)
		assert.True(t, features.FullAnalysis)
		assert.True(t, features.Projections)
		assert.True(t, features.SmartAlerts)
		assert.True(t, features.ExportData)
		assert.Equal(t, -1, features.MaxTransactions)
		assert.Len(t, features.Features, 7)
	})

	t.Run("default fallback to free", func(t *testing.T) {
		features := subscription.GetPlanFeatures("unknown")
		assert.Equal(t, "free", features.Plan)
		assert.False(t, features.FullAnalysis)
	})
}
