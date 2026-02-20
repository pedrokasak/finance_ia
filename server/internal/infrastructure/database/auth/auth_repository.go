package auth

import (
	"errors"
	"finance-ia/internal/domain/auth"
	userdomain "finance-ia/internal/domain/user"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

// FindByEmail reads from the users table (single source of truth for auth data)
func (r *AuthRepository) FindByEmail(email string) (*auth.Authentication, error) {
	var u userdomain.User
	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		return nil, gorm.ErrRecordNotFound
	}
	return userToAuth(&u), nil
}

func (r *AuthRepository) FindByID(id uuid.UUID) (*auth.Authentication, error) {
	var u userdomain.User
	if err := r.db.Where("id = ?", id).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return userToAuth(&u), nil
}

// Update persists authentication-related fields (2FA secret/enabled) back to the users table.
// This is the critical fix: auth data lives in the users table, not a separate authentications table.
func (r *AuthRepository) Update(a *auth.Authentication) error {
	return r.db.Model(&userdomain.User{}).
		Where("id = ?", a.ID).
		Updates(map[string]interface{}{
			"two_fa_enabled": a.TwoFAEnabled,
			"two_fa_secret":  a.TwoFASecret,
		}).Error
}

func (r *AuthRepository) Login(email, password string) (string, error) {
	var u userdomain.User
	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		return "", err
	}
	return u.Email, nil
}

// ResetPassword hashes and updates the user's password in the users table
func (r *AuthRepository) ResetPassword(email, newPassword string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return r.db.Model(&userdomain.User{}).
		Where("email = ?", email).
		Update("password", string(hashed)).Error
}

func (r *AuthRepository) RecoveryPassword(_ string) error {
	return nil
}

func (r *AuthRepository) ValidateToken(tokenString string) (*jwt.Token, error) {
	var jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return token, nil
}

func (r *AuthRepository) Logout(_ string) error {
	return nil
}

// userToAuth maps a User entity to the auth domain model
func userToAuth(u *userdomain.User) *auth.Authentication {
	return &auth.Authentication{
		ID:           u.ID,
		Email:        u.Email,
		Password:     u.Password,
		TwoFAEnabled: u.TwoFAEnabled,
		TwoFASecret:  u.TwoFASecret,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}
