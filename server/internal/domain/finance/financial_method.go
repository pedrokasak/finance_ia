package finance

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// FinancialMethodSplit represents the percentage allocation for a category
type FinancialMethodSplit struct {
	Label   string  `json:"label"`
	Percent float64 `json:"percent"`
	Color   string  `json:"color"`
}

// FinancialMethod represents a budgeting strategy
type FinancialMethod struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Key         string    `json:"key" gorm:"unique;not null"` // e.g., "50-30-20"
	Name        string    `json:"name" gorm:"not null"`
	Tagline     string    `json:"tagline" gorm:"not null"`
	Description string    `json:"description" gorm:"not null"`
	ForWho      string    `json:"for_who" gorm:"not null"`
	Icon        string    `json:"icon" gorm:"not null"`        // e.g., "PieChart"
	Color       string    `json:"color" gorm:"not null"`       // e.g., "text-emerald-400"
	Bg          string    `json:"bg" gorm:"not null"`          // e.g., "bg-emerald-500/10 border-emerald-500/30"
	SplitRaw    string    `json:"-" gorm:"type:text;not null"` // Store as JSON string in DB
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (m *FinancialMethod) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

// Split returns the parsed JSON array of splits
func (m *FinancialMethod) Split() []FinancialMethodSplit {
	var splits []FinancialMethodSplit
	if m.SplitRaw != "" {
		_ = json.Unmarshal([]byte(m.SplitRaw), &splits)
	}
	return splits
}

// MarshalJSON overrides default JSON marshaling to include the parsed Split array
func (m *FinancialMethod) MarshalJSON() ([]byte, error) {
	type Alias FinancialMethod
	return json.Marshal(&struct {
		*Alias
		Split []FinancialMethodSplit `json:"split"`
	}{
		Alias: (*Alias)(m),
		Split: m.Split(),
	})
}

// FinancialMethodRepository defines data access for financial methods
type FinancialMethodRepository interface {
	List() ([]*FinancialMethod, error)
	FindByID(id uuid.UUID) (*FinancialMethod, error)
	FindByKey(key string) (*FinancialMethod, error)
	Create(method *FinancialMethod) error
}
