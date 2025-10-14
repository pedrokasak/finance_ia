package email

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PasswordResetToken struct {
	ID        string    `gorm:"type:uuid;primaryKey"`
	Token     string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	UserID    string    `gorm:"type:uuid;not null;index"`
	ExpiresAt time.Time `gorm:"not null;index"`
	Used      bool      `gorm:"default:false"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

// TableName especifica o nome da tabela
func (PasswordResetToken) TableName() string {
	return "password_reset_tokens"
}

// TokenRepository gerencia tokens de reset
type TokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

func (r *TokenRepository) SaveResetToken(token *PasswordResetToken) error {
	if token.ID == "" {
		token.ID = uuid.New().String()
	}
	if token.CreatedAt.IsZero() {
		token.CreatedAt = time.Now()
	}

	result := r.db.Create(token)
	return result.Error
}

func (r *TokenRepository) FindResetToken(token string) (*PasswordResetToken, error) {
	resetToken := &PasswordResetToken{}
	
	result := r.db.Where("token = ?", token).First(resetToken)
	
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("token not found")
	}
	
	if result.Error != nil {
		return nil, result.Error
	}

	return resetToken, nil
}

func (r *TokenRepository) UpdateResetToken(token *PasswordResetToken) error {
	// ✅ Usa GORM Save ou Updates
	result := r.db.Model(&PasswordResetToken{}).
		Where("id = ?", token.ID).
		Update("used", token.Used)
	
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("token not found")
	}

	return nil
}

func (r *TokenRepository) DeleteExpiredTokens() error {
	// ✅ Usa GORM Delete com condição
	// Deleta tokens expirados há mais de 24 horas
	expiredTime := time.Now().Add(-24 * time.Hour)
	
	result := r.db.Where("expires_at < ?", expiredTime).
		Delete(&PasswordResetToken{})
	
	return result.Error
}

// FindByUserID busca tokens de um usuário específico
func (r *TokenRepository) FindByUserID(userID string) ([]PasswordResetToken, error) {
	var tokens []PasswordResetToken
	
	result := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&tokens)
	
	return tokens, result.Error
}

// InvalidateUserTokens marca todos os tokens de um usuário como usados
func (r *TokenRepository) InvalidateUserTokens(userID string) error {
	result := r.db.Model(&PasswordResetToken{}).
		Where("user_id = ? AND used = ?", userID, false).
		Update("used", true)
	
	return result.Error
}