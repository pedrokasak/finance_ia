package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	emailDomain "finance-ia/internal/domain/email"
	emailRepository "finance-ia/internal/infrastructure/database/email"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo           Repository
	emailTokenRepo *emailRepository.TokenRepository
	emailService   emailDomain.Service
}

func NewService(r Repository, et *emailRepository.TokenRepository, es emailDomain.Service) *Service {
	return &Service{repo: r, emailTokenRepo: et, emailService: es}
}

func (s *Service) Login(email, password, code string) (string, error) {
	if err := ValidateLoginFields(email, password); err != nil {
		return "", err
	}
	user, err := s.repo.FindByEmail(email)
	var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	if user.TwoFAEnabled {
		if code == "" {
			return "", errors.New("2fa_required")
		}
		if !totp.Validate(code, user.TwoFASecret) {
			return "", errors.New("invalid 2fa code")
		}
	}

	// Default plan — updated via Stripe webhook when user upgrades
	plan := "free"

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":   user.Email,
		"user_id": user.ID,
		"plan":    plan,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "message: Nao foi possivel gerar o token", err
	}
	return tokenString, nil
}

func (s *Service) Setup2FA(email string) (string, string, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", "", errors.New("user not found")
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "FinZen",
		AccountName: user.Email,
	})
	if err != nil {
		return "", "", err
	}

	user.TwoFASecret = key.Secret()
	// TwoFAEnabled stays false until verified
	if err := s.repo.Update(user); err != nil {
		return "", "", err
	}

	return key.Secret(), key.URL(), nil
}

func (s *Service) Verify2FA(email, code string) error {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return errors.New("user not found")
	}

	if user.TwoFAEnabled {
		return errors.New("2fa already enabled")
	}

	if user.TwoFASecret == "" {
		return errors.New("2fa not set up")
	}

	valid := totp.Validate(code, user.TwoFASecret)
	if !valid {
		return errors.New("invalid 2fa code")
	}

	user.TwoFAEnabled = true
	return s.repo.Update(user)
}

func (s *Service) Disable2FA(email string) error {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return errors.New("user not found")
	}

	user.TwoFAEnabled = false
	user.TwoFASecret = ""
	return s.repo.Update(user)
}

func (s *Service) GetByEmail(email string) (*Authentication, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}
	return s.repo.FindByEmail(email)
}

func (s *Service) Update(auth *Authentication) error {
	if auth.ID == uuid.Nil {
		return errors.New("id is required")
	}
	if auth.Email == "" {
		return errors.New("email is required")
	}
	return s.repo.Update(auth)
}

func (s *Service) ValidateToken(tokenString string) (*jwt.Token, error) {
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

func (s *Service) Logout(tokenString string) error {
	if tokenString == "" {
		return errors.New("token is required")
	}

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	parsed, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})
	if err != nil || !parsed.Valid {
		return errors.New("invalid token")
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid token claims")
	}

	var email string
	if e, ok := claims["email"].(string); ok && e != "" {
		email = e
	} else if uid, ok := claims["user_id"].(string); ok && uid != "" {
		_ = uid
	} else {
		return errors.New("token does not contain email/user_id")
	}
	authObj, err := s.repo.FindByEmail(email)
	if err != nil {
		return err
	}
	// Invalidate the token (this is a simple approach; consider a token blacklist for production)
	authObj.Token = ""
	if err := s.repo.Update(authObj); err != nil {
		return fmt.Errorf("failed to invalidate token: %w", err)
	}

	return nil
}

func generateSecureToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (s *Service) ForgotPassword(email string) error {
	if email == "" {
		return errors.New("email is required")
	}

	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil
	}

	// Gera token único
	token, err := generateSecureToken()
	if err != nil {
		return fmt.Errorf("failed to generate token: %w", err)
	}

	// Cria registro do token (expira em 1 hora)
	resetToken := &emailRepository.PasswordResetToken{
		Token:     token,
		UserID:    user.ID.String(),
		ExpiresAt: time.Now().Add(1 * time.Hour),
		Used:      false,
	}

	// Salva no banco
	if err := s.emailTokenRepo.SaveResetToken(resetToken); err != nil {
		return fmt.Errorf("failed to save token: %w", err)
	}

	// Envia email (em background idealmente)
	if err := s.emailService.SendPasswordReset(email, token); err != nil {
		fmt.Printf("Failed to send email to %s: %v\n", email, err)

	}

	return nil
}

func (s *Service) ResetPassword(token, newPassword string) error {
	if token == "" || newPassword == "" {
		return errors.New("token and new password are required")
	}

	// Valida comprimento da senha
	if len(newPassword) < 6 {
		return errors.New("password must be at least 6 characters")
	}

	// Busca o token
	resetToken, err := s.emailTokenRepo.FindResetToken(token)
	if err != nil {
		return errors.New("invalid or expired token")
	}

	// Verifica se já foi usado
	if resetToken.Used {
		return errors.New("token already used")
	}

	// Verifica se expirou
	if time.Now().After(resetToken.ExpiresAt) {
		return errors.New("token expired")
	}

	userID, err := uuid.Parse(resetToken.UserID)
	if err != nil {
		return errors.New("invalid user ID in token")
	}

	// Busca o usuário
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Atualiza a senha (será hasheada no repositório)
	user.Password = newPassword
	if err := s.repo.ResetPassword(user.Email, newPassword); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Marca token como usado
	resetToken.Used = true
	if err := s.emailTokenRepo.UpdateResetToken(resetToken); err != nil {
		return fmt.Errorf("failed to update token: %w", err)
	}

	return nil
}
