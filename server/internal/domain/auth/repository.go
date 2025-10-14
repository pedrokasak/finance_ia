package auth

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type Repository interface {
	FindByEmail(email string) (*Authentication, error)
	FindByID(id uuid.UUID) (*Authentication, error)
	Login(email, password string) (string, error)
	ResetPassword(email, newPassword string) error
	RecoveryPassword(email string) error
	Update(auth *Authentication) error
	ValidateToken(tokenString string) (*jwt.Token, error)
	Logout(token string) error
}
