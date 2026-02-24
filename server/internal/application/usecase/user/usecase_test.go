package user_test

import (
	"testing"

	"finance-ia/internal/application/usecase/user"
	domainUser "finance-ia/internal/domain/user"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(u *domainUser.User) error {
	args := m.Called(u)
	return args.Error(0)
}

func (m *MockUserRepository) FindAll() ([]*domainUser.User, error) {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).([]*domainUser.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) FindByID(id uuid.UUID) (*domainUser.User, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*domainUser.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) FindByEmail(email string) (*domainUser.User, error) {
	args := m.Called(email)
	if args.Get(0) != nil {
		return args.Get(0).(*domainUser.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) Update(u *domainUser.User) error {
	args := m.Called(u)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(u *domainUser.User) error {
	args := m.Called(u)
	return args.Error(0)
}

func setupUseCase() (*user.UseCase, *MockUserRepository) {
	mockRepo := new(MockUserRepository)
	svc := domainUser.NewService(mockRepo)
	uc := user.NewUseCase(svc)
	return uc, mockRepo
}

func TestUseCase_Register(t *testing.T) {
	uc, repo := setupUseCase()

	t.Run("success", func(t *testing.T) {
		repo.On("Create", mock.AnythingOfType("*user.User")).Return(nil).Once()

		u, err := uc.Register("Jane", "Doe", "jane@ex.com", "pass123")
		assert.NoError(t, err)
		assert.Equal(t, "Jane", u.FirstName)
		// O password não deve ser igual a pass123 porque o UseCase faz o Hash
		assert.NotEqual(t, "pass123", u.Password)
		repo.AssertExpectations(t)
	})

	t.Run("invalid_fields", func(t *testing.T) {
		u, err := uc.Register("", "Doe", "jane@ex.com", "pass123")
		assert.Error(t, err)
		assert.Nil(t, u)
	})
}

func TestUseCase_GetAll(t *testing.T) {
	uc, repo := setupUseCase()

	t.Run("success", func(t *testing.T) {
		users := []*domainUser.User{{FirstName: "Jane"}}
		repo.On("FindAll").Return(users, nil).Once()

		res, err := uc.GetAll()
		assert.NoError(t, err)
		assert.Len(t, res, 1)
		repo.AssertExpectations(t)
	})
}

func TestUseCase_GetByID(t *testing.T) {
	uc, repo := setupUseCase()
	id := uuid.New()

	t.Run("success", func(t *testing.T) {
		repo.On("FindByID", id).Return(&domainUser.User{ID: id}, nil).Once()
		res, err := uc.GetByID(id)
		assert.NoError(t, err)
		assert.Equal(t, id, res.ID)
		repo.AssertExpectations(t)
	})
}

func TestUseCase_GetByEmail(t *testing.T) {
	uc, repo := setupUseCase()

	t.Run("success", func(t *testing.T) {
		email := "jane@ex.com"
		repo.On("FindByEmail", email).Return(&domainUser.User{Email: email}, nil).Once()
		res, err := uc.GetByEmail(email)
		assert.NoError(t, err)
		assert.Equal(t, email, res.Email)
		repo.AssertExpectations(t)
	})
}

func TestUseCase_Update(t *testing.T) {
	uc, repo := setupUseCase()

	t.Run("success", func(t *testing.T) {
		u := &domainUser.User{ID: uuid.New(), FirstName: "Jane", Email: "jane@ex.com"}
		repo.On("Update", u).Return(nil).Once()

		err := uc.Update(u)
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})
}

func TestUseCase_Delete(t *testing.T) {
	uc, repo := setupUseCase()

	t.Run("success", func(t *testing.T) {
		u := &domainUser.User{ID: uuid.New()}
		repo.On("Delete", u).Return(nil).Once()

		err := uc.Delete(u)
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})
}
