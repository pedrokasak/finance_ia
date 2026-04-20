package finance_test

import (
	"finance-ia/internal/domain/finance"
	infraFinance "finance-ia/internal/infrastructure/database/finance"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL not set; skipping postgres integration repository tests")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	require.NoError(t, db.Exec("CREATE EXTENSION IF NOT EXISTS pgcrypto").Error)
	err = db.AutoMigrate(&finance.Transaction{}, &finance.Category{}, &finance.Budget{}, &finance.FinancialMethod{})
	require.NoError(t, err)
	require.NoError(t, db.Exec("TRUNCATE TABLE transactions, categories, budgets, financial_methods RESTART IDENTITY CASCADE").Error)

	return db
}

func TestPostgresTransactionRepository(t *testing.T) {
	db := setupTestDB(t)
	repo := infraFinance.NewTransactionRepository(db)

	userID := uuid.New()
	tx := &finance.Transaction{
		ID:          uuid.New(),
		UserID:      userID,
		Type:        finance.TransactionTypeExpense,
		Amount:      100.0,
		Description: "Groceries",
		Date:        time.Now(),
	}

	t.Run("Create", func(t *testing.T) {
		err := repo.Create(tx)
		assert.NoError(t, err)

		var saved finance.Transaction
		err = db.First(&saved, "id = ?", tx.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, tx.Description, saved.Description)
	})

	t.Run("FindByID", func(t *testing.T) {
		found, err := repo.FindByID(tx.ID)
		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, tx.Amount, found.Amount)
	})

	t.Run("FindByUser", func(t *testing.T) {
		filter := finance.TransactionFilter{UserID: userID, Limit: 10, Page: 1}
		results, total, err := repo.FindByUser(filter)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, results, 1)
	})

	t.Run("FindByPeriod", func(t *testing.T) {
		start := time.Now().Add(-24 * time.Hour)
		end := time.Now().Add(24 * time.Hour)
		results, err := repo.FindByPeriod(userID, start, end)
		assert.NoError(t, err)
		assert.Len(t, results, 1)
	})

	t.Run("Update", func(t *testing.T) {
		tx.Amount = 250.0
		err := repo.Update(tx)
		assert.NoError(t, err)

		var updated finance.Transaction
		db.First(&updated, "id = ?", tx.ID)
		assert.Equal(t, 250.0, updated.Amount)
	})

	t.Run("IdempotencyKey", func(t *testing.T) {
		tx2 := &finance.Transaction{
			ID:             uuid.New(),
			UserID:         userID,
			IdempotencyKey: "test-key-123",
		}
		repo.Create(tx2)

		found, err := repo.FindByIdempotencyKey(userID, "test-key-123")
		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, tx2.ID, found.ID)

		_, err = repo.FindByIdempotencyKey(userID, "not-exist")
		assert.Error(t, err)
	})

	t.Run("Delete", func(t *testing.T) {
		err := repo.Delete(tx.ID, userID)
		assert.NoError(t, err)

		var count int64
		db.Model(&finance.Transaction{}).Where("id = ?", tx.ID).Count(&count)
		assert.Equal(t, int64(0), count)
	})
}

func TestPostgresCategoryRepository(t *testing.T) {
	db := setupTestDB(t)
	repo := infraFinance.NewCategoryRepository(db)

	userID := uuid.New()
	cat := &finance.Category{
		ID:     uuid.New(),
		UserID: &userID,
		Name:   "Test Category",
		Type:   finance.CategoryTypeExpense,
		Color:  "#FFFFFF",
		Icon:   "icon.png",
	}

	t.Run("Create", func(t *testing.T) {
		err := repo.Create(cat)
		assert.NoError(t, err)
	})

	t.Run("FindByID", func(t *testing.T) {
		found, err := repo.FindByID(cat.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Test Category", found.Name)
	})

	t.Run("FindByUser", func(t *testing.T) {
		cats, err := repo.FindByUser(&userID)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(cats), 1)
	})

	t.Run("Update", func(t *testing.T) {
		cat.Name = "Updated Name"
		err := repo.Update(cat)
		assert.NoError(t, err)

		found, _ := repo.FindByID(cat.ID)
		assert.Equal(t, "Updated Name", found.Name)
	})

	t.Run("SeedDefaults", func(t *testing.T) {
		err := repo.SeedDefaults()
		assert.NoError(t, err)

		var count int64
		db.Model(&finance.Category{}).Where("is_default = ?", true).Count(&count)
		assert.Greater(t, count, int64(0))
	})

	t.Run("Delete", func(t *testing.T) {
		err := repo.Delete(cat.ID, userID)
		assert.NoError(t, err)

		_, err = repo.FindByID(cat.ID)
		assert.Error(t, err) // Expect not found
	})
}

func TestPostgresBudgetRepository(t *testing.T) {
	db := setupTestDB(t)
	repo := infraFinance.NewBudgetRepository(db)

	userID := uuid.New()
	b := &finance.Budget{
		ID:           uuid.New(),
		UserID:       userID,
		Period:       "2026-02",
		TotalIncome:  1000,
		NeedsPercent: 50,
	}

	t.Run("Upsert", func(t *testing.T) {
		err := repo.Upsert(b)
		assert.NoError(t, err)

		// Update
		b.TotalIncome = 2000
		err = repo.Upsert(b)
		assert.NoError(t, err)
	})

	t.Run("FindByUserAndPeriod", func(t *testing.T) {
		found, err := repo.FindByUserAndPeriod(userID, "2026-02")
		assert.NoError(t, err)
		assert.Equal(t, 2000.0, found.TotalIncome)
	})
}
