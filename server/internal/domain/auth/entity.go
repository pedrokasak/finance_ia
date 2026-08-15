package auth

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Authentication struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Token        string    `json:"token" gorm:"unique;not null"`
	Email        string    `json:"email" gorm:"unique;not null"`
	Password     string    `json:"-" gorm:"not null"`
	TwoFAEnabled bool      `json:"two_fa_enabled"`
	TwoFASecret  string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (a *Authentication) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}
