package finance_test

import (
	"finance-ia/internal/domain/finance"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAnalyzeBehavior_NoTransactions verifies that behavior analysis returns low risk on no data
func TestAnalyzeBehavior_NoTransactions(t *testing.T) {
	svc, _, _, _, _ := newService()
	userID := uuid.New()

	analysis, err := svc.AnalyzeBehavior(userID)
	require.NoError(t, err)
	assert.Equal(t, "low", analysis.RiskLevel)
}

// TestAnalyzeBehavior_WeekendSpending verifies that weekend-heavy spending is detected
func TestAnalyzeBehavior_WeekendSpending(t *testing.T) {
	svc, txRepo, _, _, _ := newService()
	userID := uuid.New()

	// Create 4 weeks of weekend transactions (Sat/Sun) with high amounts
	now := time.Now()
	for i := 0; i < 30; i++ {
		day := now.AddDate(0, 0, -i)
		wd := day.Weekday()
		amount := 10.0
		if wd == time.Saturday || wd == time.Sunday {
			amount = 500.0 // much higher on weekends
		}
		txRepo.transactions = append(txRepo.transactions, &finance.Transaction{
			UserID: userID,
			Type:   finance.TransactionTypeExpense,
			Amount: amount,
			Date:   day,
		})
	}

	analysis, err := svc.AnalyzeBehavior(userID)
	require.NoError(t, err)
	assert.Greater(t, analysis.WeekendVsWeekdayRatio, 1.0)
	assert.Greater(t, analysis.EmotionalSpendingScore, 0.0)
}

// TestAnalyzeBehavior_SalaryEffect verifies that post-salary spending spike is detected
func TestAnalyzeBehavior_SalaryEffect(t *testing.T) {
	svc, txRepo, _, _, _ := newService()
	userID := uuid.New()
	now := time.Now()

	// Salary income 20 days ago
	salaryDate := now.AddDate(0, 0, -20)
	txRepo.transactions = append(txRepo.transactions, &finance.Transaction{
		UserID: userID,
		Type:   finance.TransactionTypeIncome,
		Amount: 5000,
		Date:   salaryDate,
	})

	// High spending in 7 days after salary
	for i := 1; i <= 7; i++ {
		txRepo.transactions = append(txRepo.transactions, &finance.Transaction{
			UserID: userID,
			Type:   finance.TransactionTypeExpense,
			Amount: 300, // high daily
			Date:   salaryDate.AddDate(0, 0, i),
		})
	}

	// Normal spending afterward
	for i := 8; i <= 20; i++ {
		txRepo.transactions = append(txRepo.transactions, &finance.Transaction{
			UserID: userID,
			Type:   finance.TransactionTypeExpense,
			Amount: 50, // much lower
			Date:   salaryDate.AddDate(0, 0, i),
		})
	}

	analysis, err := svc.AnalyzeBehavior(userID)
	require.NoError(t, err)
	// Salary effect should be detected (>30% spike)
	assert.Greater(t, analysis.SalaryEffectAmount, 0.0)
	assert.Equal(t, 7, analysis.SalaryEffectDays)
}

// TestAnalyzeBehavior_ImpulsePurchases checks impulse purchase detection
func TestAnalyzeBehavior_ImpulsePurchases(t *testing.T) {
	svc, txRepo, _, _, _ := newService()
	userID := uuid.New()
	now := time.Now()

	// 15 small non-recurring purchases (< R$50)
	for i := 0; i < 15; i++ {
		txRepo.transactions = append(txRepo.transactions, &finance.Transaction{
			UserID:      userID,
			Type:        finance.TransactionTypeExpense,
			Amount:      30,
			IsRecurring: false,
			Date:        now.AddDate(0, 0, -i),
		})
	}

	analysis, err := svc.AnalyzeBehavior(userID)
	require.NoError(t, err)
	assert.Equal(t, 15, analysis.ImpulsivePurchaseCount)
	assert.InDelta(t, 450.0, analysis.ImpulsivePurchaseTotal, 0.01)
}

// TestAnalyzeBehavior_MoneyLeak verifies money leak category detection
func TestAnalyzeBehavior_MoneyLeak(t *testing.T) {
	svc, txRepo, _, _, _ := newService()
	userID := uuid.New()
	now := time.Now()

	catID := uuid.New()
	catFood := &finance.Category{ID: catID, Name: "Alimentação", Color: "#EF4444"}

	// Previous month: R$500
	prevMonth := now.AddDate(0, -1, 0)
	txRepo.transactions = append(txRepo.transactions, &finance.Transaction{
		UserID:     userID,
		Type:       finance.TransactionTypeExpense,
		Amount:     500,
		Date:       prevMonth,
		CategoryID: &catID,
		Category:   catFood,
	})

	// Current month: R$700 (40% growth)
	txRepo.transactions = append(txRepo.transactions, &finance.Transaction{
		UserID:     userID,
		Type:       finance.TransactionTypeExpense,
		Amount:     700,
		Date:       now,
		CategoryID: &catID,
		Category:   catFood,
	})

	analysis, err := svc.AnalyzeBehavior(userID)
	require.NoError(t, err)
	require.Len(t, analysis.MoneyLeakCategories, 1)
	assert.Equal(t, "Alimentação", analysis.MoneyLeakCategories[0].CategoryName)
	assert.InDelta(t, 40.0, analysis.MoneyLeakCategories[0].GrowthPercent, 0.1)
}

// TestAnalyzeBehavior_HighRisk verifies risk level aggregation
func TestAnalyzeBehavior_HighRisk(t *testing.T) {
	svc, txRepo, _, _, _ := newService()
	userID := uuid.New()
	now := time.Now()

	catID := uuid.New()
	cat := &finance.Category{ID: catID, Name: "Lazer", Color: "#EC4899"}

	// Income
	txRepo.transactions = append(txRepo.transactions, &finance.Transaction{
		UserID: userID, Type: finance.TransactionTypeIncome, Amount: 4000, Date: now.AddDate(0, -1, 0),
	})

	// High weekend spending
	for i := 0; i < 12; i++ {
		sat := now.AddDate(0, 0, -(i*7 + 1))
		txRepo.transactions = append(txRepo.transactions, &finance.Transaction{
			UserID: userID, Type: finance.TransactionTypeExpense, Amount: 400, Date: sat,
			CategoryID: &catID, Category: cat,
		})
	}
	// Normal weekdays
	for i := 0; i < 20; i++ {
		txRepo.transactions = append(txRepo.transactions, &finance.Transaction{
			UserID: userID, Type: finance.TransactionTypeExpense, Amount: 20, Date: now.AddDate(0, 0, -i),
		})
	}
	// Impulse purchases
	for i := 0; i < 12; i++ {
		txRepo.transactions = append(txRepo.transactions, &finance.Transaction{
			UserID: userID, Type: finance.TransactionTypeExpense, Amount: 25, IsRecurring: false, Date: now.AddDate(0, 0, -i),
		})
	}

	// Money leak (previous month)
	prevMonth := now.AddDate(0, -1, 0)
	txRepo.transactions = append(txRepo.transactions, &finance.Transaction{
		UserID: userID, Type: finance.TransactionTypeExpense, Amount: 200, Date: prevMonth, CategoryID: &catID, Category: cat,
	})

	analysis, err := svc.AnalyzeBehavior(userID)
	require.NoError(t, err)
	// Multiple risk factors should result in medium or high risk
	assert.NotEqual(t, "low", analysis.RiskLevel)
}
