package user

import (
	"errors"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Register(firstName, lastName, email, password string) (*User, error) {
	user := &User{FirstName: firstName, LastName: lastName, Email: email, Password: password}
	if err := ValidateUserFields(user); err != nil {
		return nil, err
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) Login(email, password string) (string, error) {
	if err := ValidateLoginFields(email, password); err != nil { return "", err }
	user, err := s.repo.FindByEmail(email)
	var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

	if err != nil {
		return "", err
	}
	if user.Password != password {
		return "error", ErrIncorrectPassword()
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

func (s *Service) GetAll() ([]*User, error) {
	users, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, errors.New("users not found")
	}
	return users, nil
}

func (s *Service) GetByID(id uuid.UUID) (*User, error) {
	if id == uuid.Nil {
        return nil, errors.New("invalid ID")
    }
	user, err := s.repo.FindByID(id)
	if err != nil {
			return nil, err
	}
  return user, nil
}

func (s *Service) GetByEmail(email string) (*User, error) {
	return s.repo.FindByEmail(email)
}

func (s *Service) Update(user *User) error {
	if user.ID == uuid.Nil {
			return errors.New("ID invalid")
	}
	if user.FirstName == "" || user.Email == "" {
			return errors.New("required fields are missing")
	}
	return s.repo.Update(user)
}

func (s *Service) Delete(user *User) error {
	if user.ID == uuid.Nil {
			return errors.New("invalid ID")
	}
	return s.repo.Delete(user)
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
	user := &User{} // Substitua isso pela lógica real de obtenção do usuário pelo token
	user.Password = newPassword
	return s.repo.Update(user)
}
