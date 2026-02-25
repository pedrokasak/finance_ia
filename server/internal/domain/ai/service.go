package ai

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// planLimits defines how many AI insights each plan gets per week
var planLimits = map[string]int{
	"free":    1,
	"premium": -1, // unlimited
	"pro":     -1, // unlimited
}

type Service struct {
	insightRepo InsightRepository
	provider    AIProvider
}

func NewService(insightRepo InsightRepository, provider AIProvider) *Service {
	return &Service{
		insightRepo: insightRepo,
		provider:    provider,
	}
}

// GetInsight returns a cached insight or generates a new one respecting plan limits
func (s *Service) GetInsight(userID uuid.UUID, plan string, ctx FinancialContext) (*AIInsight, error) {
	period := time.Now().Format("2006-01")

	// Check if there's a non-expired insight in cache
	existing, err := s.insightRepo.FindLatestByUser(userID, period)
	if err == nil && existing != nil && time.Now().Before(existing.ExpiresAt) {
		return existing, nil
	}

	// Check rate limit for free plan
	limit := planLimits[plan]
	if limit > 0 {
		weekStart := time.Now().AddDate(0, 0, -7)
		count, err := s.insightRepo.CountByUserAndPeriod(userID, weekStart, time.Now())
		if err != nil {
			return nil, err
		}
		if count >= int64(limit) {
			return nil, errors.New("rate_limited: upgrade your plan for more insights")
		}
	}

	// Generate new insight
	ctx.Plan = plan
	ctx.Period = period

	insight, err := s.provider.GenerateInsight(ctx)
	if err != nil {
		return nil, err
	}

	insight.UserID = userID
	insight.Period = period
	insight.Plan = plan
	insight.GeneratedAt = time.Now()

	// Expiration: 7 days for free, 24h for premium/pro
	if plan == "free" {
		insight.ExpiresAt = time.Now().AddDate(0, 0, 7)
	} else {
		insight.ExpiresAt = time.Now().Add(24 * time.Hour)
	}

	if err := s.insightRepo.Save(insight); err != nil {
		return nil, err
	}

	return insight, nil
}

// GetDiagnostic is available for premium and pro plans
func (s *Service) GetDiagnostic(userID uuid.UUID, plan string, ctx FinancialContext) (*AIInsight, error) {
	if plan == "free" {
		return nil, errors.New("upgrade_required: diagnostic requires Premium or Pro plan")
	}

	ctx.Plan = plan
	ctx.Period = time.Now().Format("2006-01")

	insight, err := s.provider.GenerateDiagnostic(ctx)
	if err != nil {
		return nil, err
	}

	insight.UserID = userID
	insight.Period = ctx.Period
	insight.Plan = plan
	insight.GeneratedAt = time.Now()
	insight.ExpiresAt = time.Now().Add(24 * time.Hour)
	_ = s.insightRepo.Save(insight)

	return insight, nil
}

// GetProjection is available for premium and pro plans
func (s *Service) GetProjection(userID uuid.UUID, plan string, ctx FinancialContext) (*AIInsight, error) {
	if plan == "free" {
		return nil, errors.New("upgrade_required: projection requires Premium or Pro plan")
	}

	ctx.Plan = plan
	ctx.Period = time.Now().Format("2006-01")

	insight, err := s.provider.GenerateProjection(ctx)
	if err != nil {
		return nil, err
	}

	insight.UserID = userID
	insight.Period = ctx.Period
	insight.Plan = plan
	insight.GeneratedAt = time.Now()
	insight.ExpiresAt = time.Now().Add(24 * time.Hour)
	_ = s.insightRepo.Save(insight)

	return insight, nil
}
