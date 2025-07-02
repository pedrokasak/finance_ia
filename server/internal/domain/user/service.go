package user

import "errors"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Register(firstName, lastName, email, password string) (*User, error) {
	if err := ValidateCreateUserFields(firstName, lastName, email, password); err != nil {
			return nil, err
	}

	user := &User{FirstName: firstName, LastName: lastName, Email: email, Password: password}
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) Login(email, password string) (string, error) {
	if err := ValidateLoginFields(email, password); err != nil { return "", err }
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", err
	}
	if user.Password != password {
		return "", errors.New("senha incorreta")
	}
	// Aqui você deve gerar um token JWT real
	token := "fake-jwt-token-for-" + user.Email
	return token, nil
}

func (s *Service) GetByID(id uint) (*User, error) {
	if id == 0 {
		return nil, errors.New("ID inválido")
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
	if user.ID == ([16]byte{}) {
		return errors.New("ID inválido")
	}
	if user.FirstName == "" || user.Email == "" {
		return errors.New("dados obrigatórios ausentes")
	}
	return s.repo.Update(user)
}

func (s *Service) Delete(user *User) error {
	if user.ID == ([16]byte{}) {
		return errors.New("ID inválido")
	}
	return s.repo.Delete(user)
}
