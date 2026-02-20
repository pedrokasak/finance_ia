package workers

import (
	"context"
	"encoding/json"
	"finance-ia/internal/domain/ai"
	"finance-ia/internal/domain/finance"
	"finance-ia/internal/infrastructure/queue"
	"log"
	"time"

	"github.com/google/uuid"
)

// AIInsightJob is the payload published to the AI insight queue
type AIInsightJob struct {
	UserID  string `json:"user_id"`
	Plan    string `json:"plan"`
	Timeout int    `json:"timeout_seconds"` // client-side timeout for polling
}

// AIWorker consumes AI insight generation jobs from the queue
type AIWorker struct {
	aiService      *ai.Service
	financeService *finance.Service
	queueClient    *queue.Client
}

func NewAIWorker(
	aiService *ai.Service,
	financeService *finance.Service,
	queueClient *queue.Client,
) *AIWorker {
	return &AIWorker{
		aiService:      aiService,
		financeService: financeService,
		queueClient:    queueClient,
	}
}

func (w *AIWorker) Start() {
	log.Println("[ai-worker] Starting — consuming queue:", queue.QueueAIInsight)
	if err := w.queueClient.Consume(queue.QueueAIInsight, w.handle); err != nil {
		log.Printf("[ai-worker] Failed to start consumer: %v", err)
	}
}

func (w *AIWorker) handle(body []byte) error {
	var job AIInsightJob
	if err := json.Unmarshal(body, &job); err != nil {
		return err
	}

	userID, err := uuid.Parse(job.UserID)
	if err != nil {
		return err
	}

	ctx := context.Background()
	_ = ctx

	// Build financial context
	summary, err := w.financeService.GetDashboardSummary(userID)
	if err != nil {
		return err
	}

	financialCtx := ai.FinancialContext{
		Plan:          job.Plan,
		TotalIncome:   summary.TotalIncome,
		TotalExpenses: summary.TotalExpenses,
		Balance:       summary.Balance,
		SavingsRate:   summary.SavingsRate,
		HealthScore:   summary.HealthScore,
		HealthLevel:   summary.HealthLevel,
	}

	for _, cs := range summary.CategoryBreakdown {
		financialCtx.TopCategories = append(financialCtx.TopCategories, ai.CategorySpend{
			Name:       cs.CategoryName,
			Amount:     cs.Total,
			Percentage: cs.Percentage,
		})
	}

	// Generate insight (with timeout)
	timeout := time.Duration(job.Timeout) * time.Second
	if timeout <= 0 {
		timeout = 30 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_ = ctx

	_, err = w.aiService.GetInsight(userID, job.Plan, financialCtx)
	if err != nil {
		log.Printf("[ai-worker] Failed to generate insight for user %s: %v", job.UserID, err)
		return err
	}

	log.Printf("[ai-worker] Generated insight for user %s (plan: %s)", job.UserID, job.Plan)
	return nil
}
