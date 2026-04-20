package finance

import (
	"time"

	"github.com/google/uuid"
)

// TransactionFilter holds query parameters for listing transactions
type TransactionFilter struct {
	UserID     uuid.UUID
	Type       *TransactionType
	CategoryID *uuid.UUID
	StartDate  *time.Time
	EndDate    *time.Time
	Page       int
	Limit      int
}

// TransactionRepository defines persistence operations for transactions
type TransactionRepository interface {
	Create(tx *Transaction) error
	FindByID(id uuid.UUID) (*Transaction, error)
	FindByUser(filter TransactionFilter) ([]*Transaction, int64, error)
	FindByPeriod(userID uuid.UUID, start, end time.Time) ([]*Transaction, error)
	Update(tx *Transaction) error
	Delete(id uuid.UUID, userID uuid.UUID) error
	FindByIdempotencyKey(userID uuid.UUID, key string) (*Transaction, error)
}

// CategoryRepository defines persistence operations for categories
type CategoryRepository interface {
	Create(cat *Category) error
	FindByID(id uuid.UUID) (*Category, error)
	FindByUser(userID *uuid.UUID) ([]*Category, error) // includes defaults
	Update(cat *Category) error
	Delete(id uuid.UUID, userID uuid.UUID) error
	SeedDefaults() error
}

// BudgetRepository defines persistence operations for budgets
type BudgetRepository interface {
	Upsert(budget *Budget) error
	FindByUserAndPeriod(userID uuid.UUID, period string) (*Budget, error)
}
