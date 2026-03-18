package goal_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"finance-ia/internal/domain/goal"
	goalHandler "finance-ia/internal/interfaces/handlers/goal"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
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

func setupRouter(plan string, mockRepo *MockGoalRepository) (*gin.Engine, uuid.UUID) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	userID := uuid.New()

	// Inject Plan and UserID middleware
	r.Use(func(c *gin.Context) {
		c.Set("plan", plan)
		c.Set("user_id", userID.String())
		c.Next()
	})

	svc := goal.NewService(mockRepo)
	handler := goalHandler.NewGoalHandler(svc)
	handler.RegisterRoutes(r, r)

	return r, userID
}

func TestGoalHandler_CreateGoal(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockGoalRepository)
		r, _ := setupRouter("pro", mockRepo)

		body := map[string]interface{}{
			"name":          "Casa Própria",
			"target_amount": 100000,
			"target_date":   "2030-01-01T00:00:00Z",
			"icon":          "home",
		}
		jsonValue, _ := json.Marshal(body)

		mockRepo.On("Create", mock.AnythingOfType("*goal.Goal")).Return(nil).Once()

		req, _ := http.NewRequest("POST", "/goals/", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid_plan", func(t *testing.T) {
		mockRepo := new(MockGoalRepository)
		rFree, _ := setupRouter("free", mockRepo)
		req, _ := http.NewRequest("POST", "/goals/", bytes.NewBuffer([]byte("{}")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		rFree.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("invalid_body", func(t *testing.T) {
		mockRepo := new(MockGoalRepository)
		r, _ := setupRouter("pro", mockRepo)
		req, _ := http.NewRequest("POST", "/goals/", bytes.NewBuffer([]byte("{invalid}")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service_error", func(t *testing.T) {
		mockRepo := new(MockGoalRepository)
		r, _ := setupRouter("pro", mockRepo)
		body := map[string]interface{}{
			"name":          "Carro",
			"target_amount": 50000,
			"target_date":   "2026-01-01T00:00:00Z",
		}
		jsonValue, _ := json.Marshal(body)

		mockRepo.On("Create", mock.AnythingOfType("*goal.Goal")).Return(errors.New("db error")).Once()

		req, _ := http.NewRequest("POST", "/goals/", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestGoalHandler_GetGoals(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockGoalRepository)
		r, userID := setupRouter("premium", mockRepo)

		goals := []*goal.Goal{{ID: uuid.New(), UserID: userID, Name: "Aposentadoria"}}
		mockRepo.On("FindByUserID", userID).Return(goals, nil).Once()

		req, _ := http.NewRequest("GET", "/goals/", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var res []goal.Goal
		json.Unmarshal(w.Body.Bytes(), &res)
		assert.Len(t, res, 1)
		mockRepo.AssertExpectations(t)
	})

	t.Run("service_error", func(t *testing.T) {
		mockRepo := new(MockGoalRepository)
		r, userID := setupRouter("premium", mockRepo)

		mockRepo.On("FindByUserID", userID).Return(nil, errors.New("db error")).Once()

		req, _ := http.NewRequest("GET", "/goals/", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestGoalHandler_DeleteGoal(t *testing.T) {
	goalID := uuid.New()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockGoalRepository)
		r, userID := setupRouter("pro", mockRepo)

		g := &goal.Goal{ID: goalID, UserID: userID}
		mockRepo.On("FindByID", goalID).Return(g, nil)
		mockRepo.On("Delete", g).Return(nil)

		req, _ := http.NewRequest("DELETE", "/goals/"+goalID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("forbidden_access", func(t *testing.T) {
		mockRepo := new(MockGoalRepository)
		r, _ := setupRouter("pro", mockRepo)

		g := &goal.Goal{ID: goalID, UserID: uuid.New()} // Outro dono
		mockRepo.On("FindByID", goalID).Return(g, nil)

		req, _ := http.NewRequest("DELETE", "/goals/"+goalID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}

func TestGoalHandler_UpdateGoal(t *testing.T) {
	goalID := uuid.New()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockGoalRepository)
		r, userID := setupRouter("pro", mockRepo)

		g := &goal.Goal{ID: goalID, UserID: userID, TargetAmount: 1000, CurrentAmount: 500}
		mockRepo.On("FindByID", goalID).Return(g, nil).Once()

		body := map[string]interface{}{
			"current_amount": 800,
			"target_date":    time.Now().Format(time.RFC3339),
			"icon":           "car",
		}
		jsonValue, _ := json.Marshal(body)

		mockRepo.On("Update", g).Return(nil).Once()

		req, _ := http.NewRequest("PUT", "/goals/"+goalID.String(), bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockRepo.AssertExpectations(t)
	})
}
