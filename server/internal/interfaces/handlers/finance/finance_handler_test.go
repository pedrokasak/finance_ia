package finance_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"finance-ia/internal/domain/finance"
	financeHandler "finance-ia/internal/interfaces/handlers/finance"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// ─── Simple Mocks ───────────────────────────────────────────────────────────
type mockTxRepo struct {
	transactions []*finance.Transaction
	createErr    error
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
	return m.transactions, int64(len(m.transactions)), nil
}
func (m *mockTxRepo) FindByPeriod(userID uuid.UUID, start, end time.Time) ([]*finance.Transaction, error) {
	return m.transactions, nil
}
func (m *mockTxRepo) Update(tx *finance.Transaction) error {
	for i, existing := range m.transactions {
		if existing.ID == tx.ID {
			m.transactions[i] = tx
			return nil
		}
	}
	return errors.New("not found")
}
func (m *mockTxRepo) Delete(id, userID uuid.UUID) error { return nil }
func (m *mockTxRepo) FindByIdempotencyKey(userID uuid.UUID, key string) (*finance.Transaction, error) {
	return nil, errors.New("not found")
}

type mockCategoryRepo struct {
	categories []*finance.Category
}

func (m *mockCategoryRepo) Create(cat *finance.Category) error {
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
	if budget.TotalIncome == 0 {
		return errors.New("invalid budget")
	}
	m.budget = budget
	return nil
}
func (m *mockBudgetRepo) FindByUserAndPeriod(userID uuid.UUID, period string) (*finance.Budget, error) {
	if m.budget != nil {
		return m.budget, nil
	}
	return nil, errors.New("not found")
}

type mockMethodRepo struct{}

func (m *mockMethodRepo) List() ([]*finance.FinancialMethod, error)               { return nil, nil }
func (m *mockMethodRepo) FindByID(id uuid.UUID) (*finance.FinancialMethod, error) { return nil, nil }
func (m *mockMethodRepo) FindByKey(key string) (*finance.FinancialMethod, error)  { return nil, nil }
func (m *mockMethodRepo) Create(method *finance.FinancialMethod) error            { return nil }

// ─── setup ───────────────────────────────────────────────────────────────────
func setupRouter() (*gin.Engine, *mockTxRepo, *mockCategoryRepo, *mockBudgetRepo, uuid.UUID) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	userID := uuid.New()

	r.Use(func(c *gin.Context) {
		c.Set("user_id", userID.String())
		c.Next()
	})

	txRepo := &mockTxRepo{}
	catRepo := &mockCategoryRepo{}
	budRepo := &mockBudgetRepo{}
	methRepo := &mockMethodRepo{}
	svc := finance.NewService(txRepo, catRepo, budRepo, methRepo)

	handler := financeHandler.NewFinanceHandler(svc)
	handler.RegisterRoutes(r, r)

	return r, txRepo, catRepo, budRepo, userID
}

// ─── Tests ───────────────────────────────────────────────────────────────────

func TestCreateTransactionHandler(t *testing.T) {
	r, _, _, _, _ := setupRouter()

	t.Run("success", func(t *testing.T) {
		body := map[string]interface{}{
			"type":   "expense",
			"amount": 100.50,
			"date":   "2026-02-23",
		}
		jsonValue, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/finance/transactions", bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("invalid_payload", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/finance/transactions", bytes.NewBuffer([]byte("{invalid}")))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestListTransactionsHandler(t *testing.T) {
	r, _, _, _, _ := setupRouter()

	t.Run("success", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/finance/transactions?page=1&limit=10&type=expense", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestUpdateTransactionHandler(t *testing.T) {
	r, txRepo, _, _, userID := setupRouter()

	txID := uuid.New()
	txRepo.transactions = append(txRepo.transactions, &finance.Transaction{
		ID:     txID,
		UserID: userID,
		Amount: 50.0,
		Type:   "expense",
	})

	t.Run("success", func(t *testing.T) {
		body := map[string]interface{}{
			"type":   "income",
			"amount": 200.0,
		}
		jsonValue, _ := json.Marshal(body)
		req, _ := http.NewRequest("PUT", "/finance/transactions/"+txID.String(), bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestDeleteTransactionHandler(t *testing.T) {
	r, _, _, _, _ := setupRouter()
	txID := uuid.New()

	t.Run("success", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/finance/transactions/"+txID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestCreateCategoryHandler(t *testing.T) {
	r, _, _, _, _ := setupRouter()

	t.Run("success", func(t *testing.T) {
		body := map[string]interface{}{
			"name": "Food",
			"type": "expense",
		}
		jsonValue, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/finance/categories", bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
	})
}

func TestUpsertBudgetHandler(t *testing.T) {
	r, _, _, _, _ := setupRouter()

	t.Run("success", func(t *testing.T) {
		body := map[string]interface{}{
			"total_income": 5000,
		}
		jsonValue, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/finance/budget", bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestGetDashboardHandler(t *testing.T) {
	r, _, _, _, _ := setupRouter()

	t.Run("success", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/finance/dashboard", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
