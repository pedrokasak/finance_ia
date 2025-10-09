package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	domainUser "finance-ia/internal/domain/user"
	"finance-ia/internal/interfaces/handlers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// Mock da camada de usecase
type mockUserUseCase struct{}

func (m *mockUserUseCase) Register(firstName, lastName, email, password string) (*domainUser.User, error) {
	if email == "exists@example.com" {
		return nil, errors.New("user already exists")
	}
	return &domainUser.User{
		ID:        uuid.New(),
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}, nil
}

func (m *mockUserUseCase) Login(email, password string) (string, error) {
	if email == "invalid@example.com" {
		return "", errors.New("invalid credentials")
	}
	return "mocked_token", nil
}

func (m *mockUserUseCase) GetAll() ([]*domainUser.User, error)              { return []*domainUser.User{}, nil }
func (m *mockUserUseCase) GetByID(id uuid.UUID) (*domainUser.User, error)  { return &domainUser.User{ID: id}, nil }
func (m *mockUserUseCase) GetByEmail(email string) (*domainUser.User, error) { return &domainUser.User{Email: email}, nil }
func (m *mockUserUseCase) Update(u *domainUser.User) error                 { return nil }
func (m *mockUserUseCase) Delete(u *domainUser.User) error                 { return nil }
func (m *mockUserUseCase) ForgotPassword(email string) error               { return nil }
func (m *mockUserUseCase) ResetPassword(token, newPassword string) error   { return nil }
func setupRouter(handler *handlers.UserHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/auth/signup", handler.Register)
	r.POST("/auth/login", handler.Login)
	r.POST("/auth/forgot-password", handler.ForgotPassword)
	r.POST("/auth/reset-password", handler.ResetPassword)
	return r
}

func TestRegisterUser_Success(t *testing.T) {
	usecase := &mockUserUseCase{}
	handler := handlers.NewUserHandler(usecase)
	router := setupRouter(handler)

	body := map[string]string{
		"first_name": "Pedro",
		"last_name":  "Sant Anna",
		"email":      "pedro@example.com",
		"password":   "123456",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/auth/signup", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
	assert.Contains(t, resp.Body.String(), "pedro@example.com")
}

func TestRegisterUser_AlreadyExists(t *testing.T) {
	usecase := &mockUserUseCase{}
	handler := handlers.NewUserHandler(usecase)
	router := setupRouter(handler)

	body := map[string]string{
		"first_name": "Pedro",
		"last_name":  "Sant Anna",
		"email":      "exists@example.com",
		"password":   "123456",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/auth/signup", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "user already exists")
}

func TestLogin_Success(t *testing.T) {
	usecase := &mockUserUseCase{}
	handler := handlers.NewUserHandler(usecase)
	router := setupRouter(handler)

	body := map[string]string{
		"email":    "pedro@example.com",
		"password": "123456",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "mocked_token")
}

func TestLogin_InvalidCredentials(t *testing.T) {
	usecase := &mockUserUseCase{}
	handler := handlers.NewUserHandler(usecase)
	router := setupRouter(handler)

	body := map[string]string{
		"email":    "invalid@example.com",
		"password": "wrong",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "invalid credentials")
}

func TestLogin_InvalidRequestBody(t *testing.T) {
	usecase := &mockUserUseCase{}
	handler := handlers.NewUserHandler(usecase)
	router := setupRouter(handler)

	body := "invalid json"
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "user or password is invalid")
}

func TestRecoveryPassword(t *testing.T) {
	usecase:= &mockUserUseCase{}
	handler := handlers.NewUserHandler(usecase)
	router := setupRouter(handler)

	body := map[string]string{
		"email": "pedro@example.com",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/auth/forgot-password", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "mocked_token")
}

func TestResetPassword(t *testing.T) {
	usecase:= &mockUserUseCase{}
	handler := handlers.NewUserHandler(usecase)
	router := setupRouter(handler)

	body := map[string]string{
		"token": "valid_token",
		"new_password": "newpassword123",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/auth/reset-password", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "mocked_token")
}