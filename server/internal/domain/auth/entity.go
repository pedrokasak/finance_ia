package auth

import (
	"time"

	"github.com/google/uuid"
)

type Authentication struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Token        string    `json:"token" gorm:"unique;not null"`
	Email        string    `json:"email" gorm:"unique;not null"`
	Password     string    `json:"-" gorm:"not null"`
	TwoFAEnabled bool      `json:"two_fa_enabled"`
	TwoFASecret  string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
