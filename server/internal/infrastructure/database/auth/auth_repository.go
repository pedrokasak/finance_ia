package auth

import (
	"errors"
	"finance-ia/internal/domain/auth"
	userdomain "finance-ia/internal/domain/user"
	"os"

	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) FindByEmail(email string) (*auth.Authentication, error) {
    var u userdomain.User
    if err := r.db.Where("email = ?", email).First(&u).Error; err == nil {
        return &auth.Authentication{
            ID:        u.ID,
            Email:     u.Email,
            Password:  u.Password,
            CreatedAt: u.CreatedAt,
            UpdatedAt: u.UpdatedAt,
        }, nil
    }
    return nil, gorm.ErrRecordNotFound
}

func (r *AuthRepository) Login(email, password string) (string, error) {
	var u auth.Authentication
	if err := r.db.Where("email = ? AND password = ?", email, password).First(&u).Error; err != nil {
		return "", err
	}
	return u.Token, nil
}

func (r *AuthRepository) ResetPassword(email, newPassword string) error {
	return r.db.Model(&auth.Authentication{}).Where("email = ?", email).Update("password", newPassword).Error
}

func (r *AuthRepository) RecoveryPassword(email string) error {
	var u auth.Authentication
	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		return err
	}
	// Aqui você pode gerar um token de reset e enviar por email
	_ = u // Apenas para evitar o warning de variável não utilizada
	return nil
}

func (r *AuthRepository) Update(auth *auth.Authentication) error {
	return r.db.Save(auth).Error
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