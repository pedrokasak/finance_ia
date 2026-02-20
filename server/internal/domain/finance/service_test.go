package finance_test

import (
	"errors"
	"finance-ia/internal/domain/finance"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── Mock Repositories ──────────────────────────────────────────────────────

type mockTxRepo struct {
	transactions []*finance.Transaction
	createErr    error
	findByKeyTx  *finance.Transaction
}

func (m *mockTxRepo) Create(tx *finance.Transaction) error {
	if m.createErr != nil {
		return m.createErr
	}
	if tx.ID == uuid.Nil {
		tx.ID = uuid.New()
	}
	m.transactions = append(m.transactions, tx)
	return nil
}
func (m *mockTxRepo) FindByID(id uuid.UUID) (*finance.Transaction, error) {
	for _, tx := range m.transactions {
		if tx.ID == id {
			return tx, nil
		}
	}
	return nil, errors.New("not found")
}
func (m *mockTxRepo) FindByUser(filter finance.TransactionFilter) ([]*finance.Transaction, int64, error) {
	var result []*finance.Transaction
	for _, tx := range m.transactions {
		if tx.UserID == filter.UserID {
			result = append(result, tx)
		}
	}
	return result, int64(len(result)), nil
}
func (m *mockTxRepo) FindByPeriod(userID uuid.UUID, start, end time.Time) ([]*finance.Transaction, error) {
	var result []*finance.Transaction
	for _, tx := range m.transactions {
		if tx.UserID == userID && !tx.Date.Before(start) && !tx.Date.After(end) {
			result = append(result, tx)
		}
	}
	return result, nil
}
func (m *mockTxRepo) Update(tx *finance.Transaction) error {
	for i, existing := range m.transactions {
		if existing.ID == tx.ID && existing.UserID == tx.UserID {
			m.transactions[i] = tx
			return nil
		}
	}
	return errors.New("not found")
}

func (m *mockTxRepo) Delete(id, userID uuid.UUID) error { return nil }
func (m *mockTxRepo) FindByIdempotencyKey(key string) (*finance.Transaction, error) {
	if m.findByKeyTx != nil {
		return m.findByKeyTx, nil
	}
	return nil, errors.New("not found")
}

type mockCategoryRepo struct {
	categories []*finance.Category
}

func (m *mockCategoryRepo) Create(cat *finance.Category) error {
	if cat.Name == "" {
		return errors.New("name is required")
	}
	m.categories = append(m.categories, cat)
	return nil
}
func (m *mockCategoryRepo) FindByID(id uuid.UUID) (*finance.Category, error) { return nil, nil }
func (m *mockCategoryRepo) FindByUser(_ *uuid.UUID) ([]*finance.Category, error) {
	return m.categories, nil
}
func (m *mockCategoryRepo) Update(cat *finance.Category) error { return nil }
func (m *mockCategoryRepo) Delete(id, userID uuid.UUID) error  { return nil }
func (m *mockCategoryRepo) SeedDefaults() error                { return nil }

type mockBudgetRepo struct {
	budget *finance.Budget
}

func (m *mockBudgetRepo) Upsert(budget *finance.Budget) error {
	m.budget = budget
	return nil
}
func (m *mockBudgetRepo) FindByUserAndPeriod(userID uuid.UUID, period string) (*finance.Budget, error) {
	if m.budget != nil {
		return m.budget, nil
	}
	return nil, errors.New("not found")
}

type mockMethodRepo struct {
	methods []*finance.FinancialMethod
}

func (m *mockMethodRepo) List() ([]*finance.FinancialMethod, error) {
	return m.methods, nil
}

func (m *mockMethodRepo) FindByID(id uuid.UUID) (*finance.FinancialMethod, error) {
	return nil, nil
}

func (m *mockMethodRepo) FindByKey(key string) (*finance.FinancialMethod, error) {
	return nil, nil
}

func (m *mockMethodRepo) Create(method *finance.FinancialMethod) error {
	m.methods = append(m.methods, method)
	return nil
}

// ─── Helpers ────────────────────────────────────────────────────────────────

func newService() (*finance.Service, *mockTxRepo, *mockCategoryRepo, *mockBudgetRepo, *mockMethodRepo) {
	txRepo := &mockTxRepo{}
	catRepo := &mockCategoryRepo{}
	budgetRepo := &mockBudgetRepo{}
	methodRepo := &mockMethodRepo{}
	svc := finance.NewService(txRepo, catRepo, budgetRepo, methodRepo)
	return svc, txRepo, catRepo, budgetRepo, methodRepo
}

func newTx(userID uuid.UUID, txType finance.TransactionType, amount float64) *finance.Transaction {
	return &finance.Transaction{
		UserID:      userID,
		Type:        txType,
		Amount:      amount,
		Date:        time.Now(),
		Description: "test tx",
	}
}

// ─── Transaction Tests ───────────────────────────────────────────────────────

func TestCreateTransaction_Success(t *testing.T) {
	svc, _, _, _, _ := newService()
	userID := uuid.New()

	err := svc.CreateTransaction(newTx(userID, finance.TransactionTypeExpense, 100))
	assert.NoError(t, err)
}

func TestCreateTransaction_MissingUserID(t *testing.T) {
	svc, _, _, _, _ := newService()

	err := svc.CreateTransaction(newTx(uuid.Nil, finance.TransactionTypeExpense, 100))
	assert.EqualError(t, err, "user_id is required")
}

func TestCreateTransaction_ZeroAmount(t *testing.T) {
	svc, _, _, _, _ := newService()

	err := svc.CreateTransaction(newTx(uuid.New(), finance.TransactionTypeExpense, 0))
	assert.EqualError(t, err, "amount must be greater than 0")
}

func TestCreateTransaction_NegativeAmount(t *testing.T) {
	svc, _, _, _, _ := newService()

	err := svc.CreateTransaction(newTx(uuid.New(), finance.TransactionTypeExpense, -10))
	assert.EqualError(t, err, "amount must be greater than 0")
}

func TestCreateTransaction_InvalidType(t *testing.T) {
	svc, _, _, _, _ := newService()

	tx := newTx(uuid.New(), finance.TransactionType("transfer"), 100)
	err := svc.CreateTransaction(tx)
	assert.EqualError(t, err, "type must be 'income' or 'expense'")
}

func TestCreateTransaction_IdempotencyKeyReuse(t *testing.T) {
	svc, txRepo, _, _, _ := newService()
	userID := uuid.New()
	existing := newTx(userID, finance.TransactionTypeIncome, 999)
	existing.ID = uuid.New()
	txRepo.findByKeyTx = existing

	tx := newTx(userID, finance.TransactionTypeExpense, 100)
	tx.IdempotencyKey = "duplicate-key"

	err := svc.CreateTransaction(tx)
	require.NoError(t, err)
	// Should return the existing transaction, not create a new one
	assert.Equal(t, existing.ID, tx.ID)
	assert.Equal(t, float64(999), tx.Amount)
}

func TestCreateTransaction_AutosetsDate(t *testing.T) {
	svc, _, _, _, _ := newService()
	before := time.Now()
	tx := &finance.Transaction{
		UserID: uuid.New(),
		Type:   finance.TransactionTypeIncome,
		Amount: 100,
	}
	err := svc.CreateTransaction(tx)
	require.NoError(t, err)
	assert.True(t, tx.Date.After(before) || tx.Date.Equal(before))
}

// ─── List & Delete ───────────────────────────────────────────────────────────

func TestListTransactions_DefaultLimit(t *testing.T) {
	svc, _, _, _, _ := newService()
	userID := uuid.New()

	_, _, err := svc.ListTransactions(finance.TransactionFilter{UserID: userID})
	assert.NoError(t, err)
}

func TestListTransactions_MissingUserID(t *testing.T) {
	svc, _, _, _, _ := newService()
	_, _, err := svc.ListTransactions(finance.TransactionFilter{})
	assert.EqualError(t, err, "user_id is required")
}

// ─── Budget Tests ────────────────────────────────────────────────────────────

func TestUpsertBudget_CalculatesAmounts(t *testing.T) {
	svc, _, _, budgetRepo, _ := newService()
	userID := uuid.New()

	budget := &finance.Budget{
		UserID:         userID,
		Period:         "2026-02",
		TotalIncome:    5000,
		NeedsPercent:   50,
		WantsPercent:   30,
		SavingsPercent: 20,
	}

	err := svc.UpsertBudget(budget)
	require.NoError(t, err)
	assert.Equal(t, float64(2500), budgetRepo.budget.NeedsAmount)
	assert.Equal(t, float64(1500), budgetRepo.budget.WantsAmount)
	assert.Equal(t, float64(1000), budgetRepo.budget.SavingsAmount)
}

func TestUpsertBudget_MissingUserID(t *testing.T) {
	svc, _, _, _, _ := newService()
	err := svc.UpsertBudget(&finance.Budget{TotalIncome: 1000})
	assert.EqualError(t, err, "user_id is required")
}

func TestUpsertBudget_ZeroIncome(t *testing.T) {
	svc, _, _, _, _ := newService()
	err := svc.UpsertBudget(&finance.Budget{UserID: uuid.New(), TotalIncome: 0})
	assert.EqualError(t, err, "total_income must be greater than 0")
}

// ─── Category Tests ──────────────────────────────────────────────────────────

func TestCreateCategory_EmptyName(t *testing.T) {
	svc, _, _, _, _ := newService()
	err := svc.CreateCategory(&finance.Category{})
	assert.EqualError(t, err, "name is required")
}

func TestCreateCategory_Success(t *testing.T) {
	svc, _, _, _, _ := newService()
	err := svc.CreateCategory(&finance.Category{Name: "Food", Type: finance.CategoryTypeExpense})
	assert.NoError(t, err)
}

// ─── Health Score Tests ──────────────────────────────────────────────────────

func TestGetDashboardSummary_EmptyTransactions(t *testing.T) {
	svc, _, _, _, _ := newService()
	userID := uuid.New()

	summary, err := svc.GetDashboardSummary(userID)
	require.NoError(t, err)
	assert.Equal(t, 0.0, summary.TotalIncome)
	assert.Equal(t, 0.0, summary.TotalExpenses)
	assert.Equal(t, 500, summary.HealthScore) // baseline with no transactions
}

func TestGetDashboardSummary_PositiveBalance(t *testing.T) {
	svc, txRepo, _, _, _ := newService()
	userID := uuid.New()
	now := time.Now()

	txRepo.transactions = []*finance.Transaction{
		{UserID: userID, Type: finance.TransactionTypeIncome, Amount: 5000, Date: now},
		{UserID: userID, Type: finance.TransactionTypeExpense, Amount: 2000, Date: now},
	}

	summary, err := svc.GetDashboardSummary(userID)
	require.NoError(t, err)
	assert.Equal(t, float64(5000), summary.TotalIncome)
	assert.Equal(t, float64(2000), summary.TotalExpenses)
	assert.Equal(t, float64(3000), summary.Balance)
	assert.InDelta(t, 60.0, summary.SavingsRate, 0.01)
	assert.GreaterOrEqual(t, summary.HealthScore, 700) // should be high
}

func TestGetDashboardSummary_NegativeBalance(t *testing.T) {
	svc, txRepo, _, _, _ := newService()
	userID := uuid.New()
	now := time.Now()

	txRepo.transactions = []*finance.Transaction{
		{UserID: userID, Type: finance.TransactionTypeIncome, Amount: 1000, Date: now},
		{UserID: userID, Type: finance.TransactionTypeExpense, Amount: 2000, Date: now},
	}

	summary, err := svc.GetDashboardSummary(userID)
	require.NoError(t, err)
	assert.Less(t, summary.Balance, 0.0)
	assert.Less(t, summary.HealthScore, 500) // should be below baseline
}

// ─── Health Score Level Tests ────────────────────────────────────────────────

// TestHealthScoreLevels verifies the level returned by each score range.
// We test calculateHealthScore indirectly via GetDashboardSummary.
// Health score = baseline(500) + savingsRate contribution + balance ratio contribution.
//
//	Diamante ≥900: savings 40%, balance ratio 60% → +300 +200 = 1000
//	Platina  ≥750: savings 20%, balance ratio 30% → +200 +100 = 800
//	Ouro     ≥600: savings 10%, balance ratio 15% → +100 +100 = 700
//	Prata    ≥400: savings  5%, balance ratio  5% → +0   +0   = 500  (but must not trigger -200 for negative savings)
//	Bronze    <400: negative savings, negative balance → -200 -150 = 150
func TestHealthScoreLevels(t *testing.T) {
	cases := []struct {
		name     string
		income   float64
		expense  float64
		minScore int
		maxScore int
		level    string
	}{
		// savings=40%, balance/income=0.6≥0.5 → +300+200 = 1000
		{"Diamante", 10000, 6000, 900, 1000, "Diamante"},
		// savings=20%, balance/income=0.3≥0.2 → +200+100 = 800
		{"Platina", 10000, 8000, 750, 899, "Platina"},
		// savings=10%, balance/income=0.1<0.2 → +100+0 = 600
		{"Ouro", 10000, 9000, 600, 749, "Ouro"},
		// savings=5%, balance/income=0.05<0.2 → +0+0 = 500
		{"Prata", 10000, 9500, 400, 599, "Prata"},
		// expense>income: savings=-100%, balance<0 → -200-150 = 150
		{"Bronze", 1000, 4000, 0, 399, "Bronze"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			svc, txRepo, _, _, _ := newService()
			userID := uuid.New()
			now := time.Now()

			txRepo.transactions = []*finance.Transaction{
				{UserID: userID, Type: finance.TransactionTypeIncome, Amount: tc.income, Date: now},
				{UserID: userID, Type: finance.TransactionTypeExpense, Amount: tc.expense, Date: now},
			}

			summary, err := svc.GetDashboardSummary(userID)
			require.NoError(t, err)
			assert.Equal(t, tc.level, summary.HealthLevel,
				"score=%d (expected %d-%d)", summary.HealthScore, tc.minScore, tc.maxScore)
			assert.GreaterOrEqual(t, summary.HealthScore, tc.minScore)
			assert.LessOrEqual(t, summary.HealthScore, tc.maxScore)
		})
	}
}
