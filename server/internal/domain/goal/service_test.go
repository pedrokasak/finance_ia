package goal_test

import (
	"errors"
	"testing"
	"time"

	"finance-ia/internal/domain/goal"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockGoalRepository struct {
	mock.Mock
}

func (m *MockGoalRepository) Create(g *goal.Goal) error {
	args := m.Called(g)
	return args.Error(0)
}

func (m *MockGoalRepository) FindByID(id uuid.UUID) (*goal.Goal, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*goal.Goal), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockGoalRepository) FindByUserID(userID uuid.UUID) ([]*goal.Goal, error) {
	args := m.Called(userID)
	if args.Get(0) != nil {
		return args.Get(0).([]*goal.Goal), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockGoalRepository) Update(g *goal.Goal) error {
	args := m.Called(g)
	return args.Error(0)
}

func (m *MockGoalRepository) Delete(g *goal.Goal) error {
	args := m.Called(g)
	return args.Error(0)
}

func TestGoalService_Create(t *testing.T) {
	mockRepo := new(MockGoalRepository)
	service := goal.NewService(mockRepo)

	t.Run("success", func(t *testing.T) {
		g := &goal.Goal{
			UserID:       uuid.New(),
			Name:         "Novo Carro",
			TargetAmount: 50000,
			TargetDate:   time.Now().AddDate(1, 0, 0),
		}

		mockRepo.On("Create", g).Return(nil).Once()

		err := service.Create(g)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("missing name", func(t *testing.T) {
		g := &goal.Goal{
			UserID:       uuid.New(),
			TargetAmount: 50000,
		}

		err := service.Create(g)
		assert.Error(t, err)
		assert.Equal(t, "name is required", err.Error())
	})

	t.Run("invalid target amount", func(t *testing.T) {
		g := &goal.Goal{
			UserID:       uuid.New(),
			Name:         "Novo Carro",
			TargetAmount: -10,
		}

		err := service.Create(g)
		assert.Error(t, err)
		assert.Equal(t, "target_amount must be greater than zero", err.Error())
	})
}

func TestGoalService_Update(t *testing.T) {
	mockRepo := new(MockGoalRepository)
	service := goal.NewService(mockRepo)

	t.Run("success", func(t *testing.T) {
		g := &goal.Goal{
			ID:            uuid.New(),
			CurrentAmount: 100,
		}

		mockRepo.On("Update", g).Return(nil).Once()

		err := service.Update(g)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("negative current amount", func(t *testing.T) {
		g := &goal.Goal{
			ID:            uuid.New(),
			CurrentAmount: -50,
		}

		err := service.Update(g)
		assert.Error(t, err)
		assert.Equal(t, "current_amount cannot be negative", err.Error())
	})
}

func TestGoalService_Delete(t *testing.T) {
	mockRepo := new(MockGoalRepository)
	service := goal.NewService(mockRepo)

	t.Run("success", func(t *testing.T) {
		goalID := uuid.New()
		g := &goal.Goal{ID: goalID}

		mockRepo.On("FindByID", goalID).Return(g, nil).Once()
		mockRepo.On("Delete", g).Return(nil).Once()

		err := service.Delete(goalID)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		goalID := uuid.New()
		mockRepo.On("FindByID", goalID).Return(nil, errors.New("not found")).Once()

		err := service.Delete(goalID)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
