package goal

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Goal struct {
	ID            uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	UserID        uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	Name          string    `json:"name" gorm:"not null"`
	TargetAmount  float64   `json:"target_amount" gorm:"not null"`
	CurrentAmount float64   `json:"current_amount" gorm:"default:0"`
	TargetDate    time.Time `json:"target_date" gorm:"not null"`
	Icon          string    `json:"icon" gorm:"default:'flag'"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (g *Goal) BeforeCreate(tx *gorm.DB) error {
	if g.ID == uuid.Nil {
		g.ID = uuid.New()
	}
	return nil
}
