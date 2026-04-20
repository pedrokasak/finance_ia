package finance

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	txRepo       TransactionRepository
	categoryRepo CategoryRepository
	budgetRepo   BudgetRepository
	methodRepo   FinancialMethodRepository
}

func NewService(txRepo TransactionRepository, categoryRepo CategoryRepository, budgetRepo BudgetRepository, methodRepo FinancialMethodRepository) *Service {
	return &Service{
		txRepo:       txRepo,
		categoryRepo: categoryRepo,
		budgetRepo:   budgetRepo,
		methodRepo:   methodRepo,
	}
}

// --- Transactions ---

func (s *Service) CreateTransaction(tx *Transaction) error {
	if tx.UserID == uuid.Nil {
		return errors.New("user_id is required")
	}
	if tx.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}
	if tx.Type != TransactionTypeIncome && tx.Type != TransactionTypeExpense {
		return errors.New("type must be 'income' or 'expense'")
	}
	if tx.Date.IsZero() {
		tx.Date = time.Now()
	}

	// Check idempotency
	if tx.IdempotencyKey != "" {
		existing, err := s.txRepo.FindByIdempotencyKey(tx.UserID, tx.IdempotencyKey)
		if err == nil && existing != nil {
			*tx = *existing
			return nil // Already created — idempotent response
		}
	}

	return s.txRepo.Create(tx)
}

func (s *Service) ListTransactions(filter TransactionFilter) ([]*Transaction, int64, error) {
	if filter.UserID == uuid.Nil {
		return nil, 0, errors.New("user_id is required")
	}
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	return s.txRepo.FindByUser(filter)
}

func (s *Service) UpdateTransaction(id uuid.UUID, userID uuid.UUID, payload *Transaction) (*Transaction, error) {
	existing, err := s.txRepo.FindByID(id)
	if err != nil || existing == nil {
		return nil, errors.New("transaction not found")
	}
	if existing.UserID != userID {
		return nil, errors.New("unauthorized to update this transaction")
	}

	if payload.Amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}
	if payload.Type != TransactionTypeIncome && payload.Type != TransactionTypeExpense {
		return nil, errors.New("type must be 'income' or 'expense'")
	}

	existing.Amount = payload.Amount
	existing.Type = payload.Type
	existing.Description = payload.Description
	existing.CategoryID = payload.CategoryID
	existing.IsRecurring = payload.IsRecurring
	if !payload.Date.IsZero() {
		existing.Date = payload.Date
	}

	err = s.txRepo.Update(existing)
	if err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *Service) DeleteTransaction(id uuid.UUID, userID uuid.UUID) error {
	return s.txRepo.Delete(id, userID)
}

// --- Categories ---

func (s *Service) GetCategories(userID uuid.UUID) ([]*Category, error) {
	return s.categoryRepo.FindByUser(&userID)
}

func (s *Service) GetDefaultCategories() ([]*Category, error) {
	return s.categoryRepo.FindByUser(nil)
}

func (s *Service) CreateCategory(cat *Category) error {
	if cat.Name == "" {
		return errors.New("name is required")
	}
	return s.categoryRepo.Create(cat)
}

func (s *Service) UpdateCategory(id uuid.UUID, userID uuid.UUID, payload *Category) (*Category, error) {
	existing, err := s.categoryRepo.FindByID(id)
	if err != nil || existing == nil {
		return nil, errors.New("category not found")
	}
	if existing.IsDefault {
		return nil, errors.New("cannot update default categories")
	}
	if existing.UserID == nil || *existing.UserID != userID {
		return nil, errors.New("unauthorized to update this category")
	}

	if payload.Name != "" {
		existing.Name = payload.Name
	}
	if payload.Color != "" {
		existing.Color = payload.Color
	}
	if payload.Icon != "" {
		existing.Icon = payload.Icon
	}

	if err := s.categoryRepo.Update(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *Service) DeleteCategory(id uuid.UUID, userID uuid.UUID) error {
	return s.categoryRepo.Delete(id, userID)
}

// --- Budget ---

func (s *Service) GetFinancialMethods() ([]*FinancialMethod, error) {
	if s.methodRepo == nil {
		return nil, errors.New("financial method repository is not configured")
	}
	return s.methodRepo.List()
}

func (s *Service) UpsertBudget(budget *Budget) error {
	if budget.UserID == uuid.Nil {
		return errors.New("user_id is required")
	}
	if budget.TotalIncome <= 0 {
		return errors.New("total_income must be greater than 0")
	}
	// Calculate amounts from percentages
	budget.NeedsAmount = budget.TotalIncome * budget.NeedsPercent / 100
	budget.WantsAmount = budget.TotalIncome * budget.WantsPercent / 100
	budget.SavingsAmount = budget.TotalIncome * budget.SavingsPercent / 100
	return s.budgetRepo.Upsert(budget)
}

func (s *Service) GetCurrentBudget(userID uuid.UUID) (*Budget, error) {
	period := time.Now().Format("2006-01")
	return s.budgetRepo.FindByUserAndPeriod(userID, period)
}

// --- Dashboard ---

func (s *Service) GetDashboardSummary(userID uuid.UUID) (*DashboardSummary, error) {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Second)

	transactions, err := s.txRepo.FindByPeriod(userID, startOfMonth, endOfMonth)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transactions: %w", err)
	}

	// Budget
	budget, _ := s.GetCurrentBudget(userID)

	var totalIncome, totalExpenses float64
	if budget != nil {
		totalIncome += budget.TotalIncome
	}
	categoryTotals := make(map[uuid.UUID]*CategorySummary)

	for _, tx := range transactions {
		if tx.Type == TransactionTypeIncome {
			totalIncome += tx.Amount
		} else {
			totalExpenses += tx.Amount
			if tx.CategoryID != nil && tx.Category != nil {
				if _, ok := categoryTotals[*tx.CategoryID]; !ok {
					categoryTotals[*tx.CategoryID] = &CategorySummary{
						CategoryID:   *tx.CategoryID,
						CategoryName: tx.Category.Name,
						Color:        tx.Category.Color,
					}
				}
				categoryTotals[*tx.CategoryID].Total += tx.Amount
			}
		}
	}

	balance := totalIncome - totalExpenses

	// Calculate savings rate
	savingsRate := 0.0
	if totalIncome > 0 {
		savingsRate = (balance / totalIncome) * 100
	}

	// Build category breakdown with percentages
	breakdown := make([]CategorySummary, 0, len(categoryTotals))
	for _, cs := range categoryTotals {
		if totalExpenses > 0 {
			cs.Percentage = (cs.Total / totalExpenses) * 100
		}
		breakdown = append(breakdown, *cs)
	}

	// Health score calculation (0-1000)
	healthScore := calculateHealthScore(savingsRate, balance, totalIncome)
	healthLevel := healthScoreToLevel(healthScore)

	// Days until negative (simple linear projection)
	var daysUntilNegative *int
	if balance < 0 {
		days := 0
		daysUntilNegative = &days
	} else if totalExpenses > 0 {
		dailyBurn := totalExpenses / float64(now.Day())
		if dailyBurn > 0 && totalIncome > 0 {
			days := int(balance / dailyBurn)
			if days < 30 {
				daysUntilNegative = &days
			}
		}
	}

	// Monthly trend (last 6 months)
	trend, _ := s.buildMonthlyTrend(userID, 6)

	return &DashboardSummary{
		TotalIncome:       totalIncome,
		TotalExpenses:     totalExpenses,
		Balance:           balance,
		SavingsRate:       savingsRate,
		HealthScore:       healthScore,
		HealthLevel:       healthLevel,
		Budget:            budget,
		CategoryBreakdown: breakdown,
		MonthlyTrend:      trend,
		DaysUntilNegative: daysUntilNegative,
	}, nil
}

func (s *Service) buildMonthlyTrend(userID uuid.UUID, months int) ([]MonthlyTrend, error) {
	now := time.Now()
	trends := make([]MonthlyTrend, 0, months)

	for i := months - 1; i >= 0; i-- {
		t := now.AddDate(0, -i, 0)
		start := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.Local)
		end := start.AddDate(0, 1, 0).Add(-time.Second)

		txs, err := s.txRepo.FindByPeriod(userID, start, end)
		if err != nil {
			continue
		}

		var income, expenses float64

		// Consider base budget if available for this specific month (defaulting to current if unversioned)
		periodStr := start.Format("2006-01")
		budget, _ := s.budgetRepo.FindByUserAndPeriod(userID, periodStr)
		if budget == nil {
			budget, _ = s.GetCurrentBudget(userID) // fallback to current if history isn't strict
		}
		if budget != nil {
			income += budget.TotalIncome
		}

		for _, tx := range txs {
			if tx.Type == TransactionTypeIncome {
				income += tx.Amount
			} else {
				expenses += tx.Amount
			}
		}
		trends = append(trends, MonthlyTrend{
			Month:    t.Format("Jan"),
			Income:   income,
			Expenses: expenses,
		})
	}
	return trends, nil
}

func calculateHealthScore(savingsRate, balance, income float64) int {
	score := 0 // baseline

	// Savings rate contribution (max +500)
	switch {
	case savingsRate >= 30:
		score += 500
	case savingsRate >= 20:
		score += 350
	case savingsRate >= 10:
		score += 200
	case savingsRate >= 0:
		score += 50
	default: // negative
		score += 0
	}

	// Balance contribution (max +500)
	if income > 0 {
		balanceRatio := balance / income
		if balanceRatio >= 0.5 {
			score += 500
		} else if balanceRatio >= 0.2 {
			score += 250
		} else if balanceRatio < 0 {
			score += 0
		}
	}

	if score < 0 {
		score = 0
	}
	if score > 1000 {
		score = 1000
	}
	return score
}

func healthScoreToLevel(score int) string {
	switch {
	case score >= 900:
		return "Diamante"
	case score >= 750:
		return "Platina"
	case score >= 600:
		return "Ouro"
	case score >= 400:
		return "Prata"
	default:
		return "Bronze"
	}
}
