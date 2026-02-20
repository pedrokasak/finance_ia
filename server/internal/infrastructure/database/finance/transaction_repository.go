package finance

import (
	"finance-ia/internal/domain/finance"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(tx *finance.Transaction) error {
	return r.db.Create(tx).Error
}

func (r *TransactionRepository) FindByID(id uuid.UUID) (*finance.Transaction, error) {
	var tx finance.Transaction
	if err := r.db.Preload("Category").First(&tx, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &tx, nil
}

func (r *TransactionRepository) FindByUser(filter finance.TransactionFilter) ([]*finance.Transaction, int64, error) {
	query := r.db.Model(&finance.Transaction{}).
		Preload("Category").
		Where("user_id = ?", filter.UserID)

	if filter.Type != nil {
		query = query.Where("type = ?", *filter.Type)
	}
	if filter.CategoryID != nil {
		query = query.Where("category_id = ?", *filter.CategoryID)
	}
	if filter.StartDate != nil {
		query = query.Where("date >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("date <= ?", *filter.EndDate)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.Limit
	if offset < 0 {
		offset = 0
	}

	var txs []*finance.Transaction
	if err := query.
		Order("date DESC").
		Limit(filter.Limit).
		Offset(offset).
		Find(&txs).Error; err != nil {
		return nil, 0, err
	}

	return txs, total, nil
}

func (r *TransactionRepository) FindByPeriod(userID uuid.UUID, start, end time.Time) ([]*finance.Transaction, error) {
	var txs []*finance.Transaction
	if err := r.db.
		Preload("Category").
		Where("user_id = ? AND date BETWEEN ? AND ?", userID, start, end).
		Find(&txs).Error; err != nil {
		return nil, fmt.Errorf("find transactions by period: %w", err)
	}
	return txs, nil
}

func (r *TransactionRepository) Update(tx *finance.Transaction) error {
	return r.db.Save(tx).Error
}

func (r *TransactionRepository) Delete(id uuid.UUID, userID uuid.UUID) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&finance.Transaction{}).Error
}

func (r *TransactionRepository) FindByIdempotencyKey(key string) (*finance.Transaction, error) {
	var tx finance.Transaction
	if err := r.db.Where("idempotency_key = ?", key).First(&tx).Error; err != nil {
		return nil, err
	}
	return &tx, nil
}
