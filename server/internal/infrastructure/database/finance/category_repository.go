package finance

import (
	"finance-ia/internal/domain/finance"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(cat *finance.Category) error {
	return r.db.Create(cat).Error
}

func (r *CategoryRepository) FindByID(id uuid.UUID) (*finance.Category, error) {
	var cat finance.Category
	if err := r.db.First(&cat, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &cat, nil
}

// FindByUser returns categories for a specific user + all default categories
func (r *CategoryRepository) FindByUser(userID *uuid.UUID) ([]*finance.Category, error) {
	var cats []*finance.Category
	query := r.db.Where("is_default = ?", true)
	if userID != nil {
		query = query.Or("user_id = ?", *userID)
	}
	if err := query.Order("is_default DESC, name ASC").Find(&cats).Error; err != nil {
		return nil, err
	}
	return cats, nil
}

func (r *CategoryRepository) Update(cat *finance.Category) error {
	return r.db.Save(cat).Error
}

func (r *CategoryRepository) Delete(id uuid.UUID, userID uuid.UUID) error {
	return r.db.Where("id = ? AND user_id = ? AND is_default = false", id, userID).
		Delete(&finance.Category{}).Error
}

// SeedDefaults inserts the default categories if they don't exist yet
func (r *CategoryRepository) SeedDefaults() error {
	defaults := []finance.Category{
		// Expenses
		{Name: "Alimentação", Type: finance.CategoryTypeExpense, Color: "#EF4444", Icon: "utensils", IsDefault: true},
		{Name: "Transporte", Type: finance.CategoryTypeExpense, Color: "#F59E0B", Icon: "car", IsDefault: true},
		{Name: "Saúde", Type: finance.CategoryTypeExpense, Color: "#10B981", Icon: "heart-pulse", IsDefault: true},
		{Name: "Moradia", Type: finance.CategoryTypeExpense, Color: "#3B82F6", Icon: "home", IsDefault: true},
		{Name: "Educação", Type: finance.CategoryTypeExpense, Color: "#8B5CF6", Icon: "graduation-cap", IsDefault: true},
		{Name: "Lazer", Type: finance.CategoryTypeExpense, Color: "#EC4899", Icon: "gamepad-2", IsDefault: true},
		{Name: "Roupas", Type: finance.CategoryTypeExpense, Color: "#06B6D4", Icon: "shirt", IsDefault: true},
		{Name: "Investimentos", Type: finance.CategoryTypeExpense, Color: "#84CC16", Icon: "trending-up", IsDefault: true},
		{Name: "Outros", Type: finance.CategoryTypeExpense, Color: "#6B7280", Icon: "more-horizontal", IsDefault: true},
		// Income
		{Name: "Salário", Type: finance.CategoryTypeIncome, Color: "#10B981", Icon: "briefcase", IsDefault: true},
		{Name: "Freelance", Type: finance.CategoryTypeIncome, Color: "#3B82F6", Icon: "laptop", IsDefault: true},
		{Name: "Investimentos", Type: finance.CategoryTypeIncome, Color: "#8B5CF6", Icon: "coins", IsDefault: true},
		{Name: "Outros", Type: finance.CategoryTypeIncome, Color: "#6B7280", Icon: "plus-circle", IsDefault: true},
	}

	for i := range defaults {
		r.db.Where("name = ? AND type = ? AND is_default = true", defaults[i].Name, defaults[i].Type).
			FirstOrCreate(&defaults[i])
	}
	return nil
}
