package user_test

import (
	"errors"
	"finance-ia/internal/domain/user"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(u *user.User) error {
	args := m.Called(u)
	return args.Error(0)
}

func (m *MockUserRepository) FindAll() ([]*user.User, error) {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).([]*user.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) FindByID(id uuid.UUID) (*user.User, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*user.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) FindByEmail(email string) (*user.User, error) {
	args := m.Called(email)
	if args.Get(0) != nil {
		return args.Get(0).(*user.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) Update(u *user.User) error {
	args := m.Called(u)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(u *user.User) error {
	args := m.Called(u)
	return args.Error(0)
}

func TestUserService_Register(t *testing.T) {
	repo := new(MockUserRepository)
	svc := user.NewService(repo)

	t.Run("success", func(t *testing.T) {
		repo.On("Create", mock.AnythingOfType("*user.User")).Return(nil).Once()

		u, err := svc.Register("John", "Doe", "john@example.com", "secure123")
		assert.NoError(t, err)
		assert.NotNil(t, u)
		assert.Equal(t, "John", u.FirstName)
		repo.AssertExpectations(t)
	})

	t.Run("invalid_fields", func(t *testing.T) {
		u, err := svc.Register("", "Doe", "invalid", "123")
		assert.Error(t, err)
		assert.Nil(t, u)
	})

	t.Run("repo_error", func(t *testing.T) {
		repo.On("Create", mock.AnythingOfType("*user.User")).Return(errors.New("db error")).Once()
		u, err := svc.Register("John", "Doe", "john2@example.com", "secure123")
		assert.Error(t, err)
		assert.Nil(t, u)
		repo.AssertExpectations(t)
	})
}

func TestUserService_GetAll(t *testing.T) {
	repo := new(MockUserRepository)
	svc := user.NewService(repo)

	t.Run("success", func(t *testing.T) {
		users := []*user.User{{FirstName: "John"}}
		repo.On("FindAll").Return(users, nil).Once()

		res, err := svc.GetAll()
		assert.NoError(t, err)
		assert.Len(t, res, 1)
		repo.AssertExpectations(t)
	})

	t.Run("empty", func(t *testing.T) {
		repo.On("FindAll").Return([]*user.User{}, nil).Once()
		res, err := svc.GetAll()
		assert.Error(t, err)
		assert.Equal(t, "users not found", err.Error())
		assert.Nil(t, res)
		repo.AssertExpectations(t)
	})

	t.Run("repo_error", func(t *testing.T) {
		repo.On("FindAll").Return(nil, errors.New("db error")).Once()
		res, err := svc.GetAll()
		assert.Error(t, err)
		assert.Nil(t, res)
		repo.AssertExpectations(t)
	})
}

func TestUserService_GetByID(t *testing.T) {
	repo := new(MockUserRepository)
	svc := user.NewService(repo)
	id := uuid.New()

	t.Run("success", func(t *testing.T) {
		repo.On("FindByID", id).Return(&user.User{ID: id}, nil).Once()
		res, err := svc.GetByID(id)
		assert.NoError(t, err)
		assert.Equal(t, id, res.ID)
		repo.AssertExpectations(t)
	})

	t.Run("invalid_id", func(t *testing.T) {
		res, err := svc.GetByID(uuid.Nil)
		assert.Error(t, err)
		assert.Equal(t, "invalid ID", err.Error())
		assert.Nil(t, res)
	})

	t.Run("not_found", func(t *testing.T) {
		repo.On("FindByID", id).Return(nil, errors.New("not found")).Once()
		res, err := svc.GetByID(id)
		assert.Error(t, err)
		assert.Nil(t, res)
		repo.AssertExpectations(t)
	})
}

func TestUserService_GetByEmail(t *testing.T) {
	repo := new(MockUserRepository)
	svc := user.NewService(repo)

	t.Run("success", func(t *testing.T) {
		email := "test@test.com"
		repo.On("FindByEmail", email).Return(&user.User{Email: email}, nil).Once()

		res, err := svc.GetByEmail(email)
		assert.NoError(t, err)
		assert.Equal(t, email, res.Email)
		repo.AssertExpectations(t)
	})
}

func TestUserService_Update(t *testing.T) {
	repo := new(MockUserRepository)
	svc := user.NewService(repo)

	t.Run("success", func(t *testing.T) {
		u := &user.User{ID: uuid.New(), FirstName: "John", Email: "john@ex.com"}
		repo.On("Update", u).Return(nil).Once()

		err := svc.Update(u)
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("invalid_id", func(t *testing.T) {
		err := svc.Update(&user.User{ID: uuid.Nil})
		assert.Error(t, err)
		assert.Equal(t, "ID invalid", err.Error())
	})

	t.Run("missing_fields", func(t *testing.T) {
		err := svc.Update(&user.User{ID: uuid.New(), FirstName: ""})
		assert.Error(t, err)
		assert.Equal(t, "required fields are missing", err.Error())
	})
}

func TestUserService_Delete(t *testing.T) {
	repo := new(MockUserRepository)
	svc := user.NewService(repo)

	t.Run("success", func(t *testing.T) {
		u := &user.User{ID: uuid.New()}
		repo.On("Delete", u).Return(nil).Once()

		err := svc.Delete(u)
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("invalid_id", func(t *testing.T) {
		err := svc.Delete(&user.User{ID: uuid.Nil})
		assert.Error(t, err)
		assert.Equal(t, "invalid ID", err.Error())
	})
}
