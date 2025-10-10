package auth

import (
	"errors"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo     Repository
}

func NewService(r Repository) *Service {
	return &Service{repo: r}
}


func (s *Service) Login(email, password string) (string, error) {
	if err := ValidateLoginFields(email, password); err != nil { return "", err }
	user, err := s.repo.FindByEmail(email)
	var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"user_id": user.ID,
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "message: Nao foi possivel gerar o token", err
	}
	return tokenString, nil
}

func (s *Service) ForgotPassword(email string) error {
	if email == "" {
		return errors.New("email is required")
	}
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return err
	}

	// Aqui você pode gerar um token de reset e enviar por email
	_ = user // Apenas para evitar o warning de variável não utilizada
	return nil
}

func (s *Service) ResetPassword(token, newPassword string) error {
	if token == "" || newPassword == "" {
		return errors.New("token and new password are required")
	}
	
	// Aqui você deve validar o token e encontrar o usuário correspondente
	// Para simplificação, vamos assumir que o token é válido e corresponde a um usuário
	user := &Authentication{} // Substitua isso pela lógica real de obtenção do usuário pelo token
	user.Password = newPassword
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