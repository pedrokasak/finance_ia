package ai_test

import (
	"errors"
	"finance-ia/internal/domain/ai"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── Mocks ───────────────────────────────────────────────────────────────────

type mockInsightRepo struct {
	saved  []*ai.AIInsight
	count  int64
	latest *ai.AIInsight
}

func (m *mockInsightRepo) Save(insight *ai.AIInsight) error {
	m.saved = append(m.saved, insight)
	return nil
}

func (m *mockInsightRepo) FindLatestByUser(userID uuid.UUID, period string) (*ai.AIInsight, error) {
	if m.latest != nil {
		return m.latest, nil
	}
	return nil, errors.New("not found")
}

func (m *mockInsightRepo) CountByUserAndPeriod(userID uuid.UUID, start, end time.Time) (int64, error) {
	return m.count, nil
}

type mockAIProvider struct {
	insight  *ai.AIInsight
	insights []*ai.AIInsight
	err      error
}

func (m *mockAIProvider) GenerateInsight(ctx ai.FinancialContext) (*ai.AIInsight, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.insight != nil {
		return m.insight, nil
	}
	return &ai.AIInsight{
		ID:          uuid.New(),
		Type:        ai.InsightTypeTip,
		Title:       "Mock Insight",
		Content:     "Save more money.",
		GeneratedAt: time.Now(),
		ExpiresAt:   time.Now().Add(7 * 24 * time.Hour),
		Period:      time.Now().Format("2006-01"),
	}, nil
}

func (m *mockAIProvider) GenerateFullAnalysis(ctx ai.FinancialContext) ([]*ai.AIInsight, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.insights != nil {
		return m.insights, nil
	}
	return []*ai.AIInsight{
		{ID: uuid.New(), Type: ai.InsightTypeAchievement, Title: "Full Analysis", Content: "You're doing great."},
	}, nil
}

func newSvc(repo *mockInsightRepo, provider ai.AIProvider) *ai.Service {
	return ai.NewService(repo, provider)
}

func baseCtx(plan string) ai.FinancialContext {
	return ai.FinancialContext{
		Plan:          plan,
		TotalIncome:   5000,
		TotalExpenses: 3000,
		Balance:       2000,
		SavingsRate:   40,
		HealthScore:   750,
		HealthLevel:   "Platina",
	}
}

// ─── Free Plan Tests ──────────────────────────────────────────────────────────

func TestGetInsight_FreePlan_GeneratesFirstTime(t *testing.T) {
	repo := &mockInsightRepo{count: 0}
	svc := newSvc(repo, &mockAIProvider{})
	userID := uuid.New()

	insight, err := svc.GetInsight(userID, "free", baseCtx("free"))
	require.NoError(t, err)
	assert.NotNil(t, insight)
	assert.Len(t, repo.saved, 1)
}

func TestGetInsight_FreePlan_CacheHitReturnsExisting(t *testing.T) {
	existing := &ai.AIInsight{
		ID:        uuid.New(),
		Title:     "Cached Insight",
		Period:    time.Now().Format("2006-01"),
		ExpiresAt: time.Now().Add(time.Hour),
	}
	repo := &mockInsightRepo{count: 0, latest: existing}
	provider := &mockAIProvider{}
	svc := newSvc(repo, provider)

	insight, err := svc.GetInsight(uuid.New(), "free", baseCtx("free"))
	require.NoError(t, err)
	assert.Equal(t, existing.ID, insight.ID)
	// Provider should NOT be called
	assert.Len(t, repo.saved, 0)
}

func TestGetInsight_FreePlan_RateLimitBlocks(t *testing.T) {
	repo := &mockInsightRepo{count: 1, latest: nil} // already used limit
	svc := newSvc(repo, &mockAIProvider{})

	_, err := svc.GetInsight(uuid.New(), "free", baseCtx("free"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "rate_limited")
}

func TestGetInsight_FreePlan_ProviderError(t *testing.T) {
	repo := &mockInsightRepo{count: 0}
	provider := &mockAIProvider{err: errors.New("gemini unavailable")}
	svc := newSvc(repo, provider)

	_, err := svc.GetInsight(uuid.New(), "free", baseCtx("free"))
	require.Error(t, err)
}

// ─── Premium Plan Tests ───────────────────────────────────────────────────────

func TestGetInsight_PremiumPlan_NoRateLimit(t *testing.T) {
	// Premium should NOT rate limit even with count=100
	repo := &mockInsightRepo{count: 100}
	svc := newSvc(repo, &mockAIProvider{})

	insight, err := svc.GetInsight(uuid.New(), "premium", baseCtx("premium"))
	require.NoError(t, err)
	assert.NotNil(t, insight)
}

func TestGetInsight_ProPlan_NoRateLimit(t *testing.T) {
	repo := &mockInsightRepo{count: 100}
	svc := newSvc(repo, &mockAIProvider{})

	insight, err := svc.GetInsight(uuid.New(), "pro", baseCtx("pro"))
	require.NoError(t, err)
	assert.NotNil(t, insight)
}

// ─── Full Analysis Tests ──────────────────────────────────────────────────────

func TestGetFullAnalysis_FreePlan_Blocked(t *testing.T) {
	repo := &mockInsightRepo{}
	svc := newSvc(repo, &mockAIProvider{})

	_, err := svc.GetFullAnalysis(uuid.New(), "free", baseCtx("free"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "upgrade_required")
}

func TestGetFullAnalysis_PremiumPlan_Allowed(t *testing.T) {
	repo := &mockInsightRepo{}
	svc := newSvc(repo, &mockAIProvider{})

	insights, err := svc.GetFullAnalysis(uuid.New(), "premium", baseCtx("premium"))
	require.NoError(t, err)
	assert.NotEmpty(t, insights)
}

func TestGetFullAnalysis_ProPlan_Allowed(t *testing.T) {
	repo := &mockInsightRepo{}
	svc := newSvc(repo, &mockAIProvider{})

	insights, err := svc.GetFullAnalysis(uuid.New(), "pro", baseCtx("pro"))
	require.NoError(t, err)
	assert.NotEmpty(t, insights)
}
