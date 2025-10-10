package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	domainUser "finance-ia/internal/domain/user"
	handlers "finance-ia/internal/interfaces/handlers/user"
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
	r.POST("/user/register", handler.Register)
	r.PATCH("/user/update/:id", handler.UpdateUser)
	r.DELETE("/user/delete/:id", handler.DeleteUser)
	r.GET("/user/:id", handler.GetUserByID)
	r.GET("/users", handler.GetAllUsers)
	return r
}

func TestRegisterUser_Success(t *testing.T) {
	usecase := &mockUserUseCase{}
	handler := handlers.NewUserHandler(usecase)
	router := setupRouter(handler)

	body := map[string]string{
		"first_name": "Pedro",
		"last_name":  "Sant Anna",
		"email":      "pedro@test.com",
		"password":   "123456",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/user/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
	assert.Contains(t, resp.Body.String(), "pedro@test.com")
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

	req, _ := http.NewRequest("POST", "/user/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "email already in use")
}

func TestGetUserByID_Success(t *testing.T) {
	usecase := &mockUserUseCase{}
	handler := handlers.NewUserHandler(usecase)
	router := setupRouter(handler)

	userID := uuid.New()
	req, _ := http.NewRequest("GET", "/user/"+userID.String(), nil)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), userID.String())
}

func TestGetUserByID_NotFound(t *testing.T) {
	usecase := &mockUserUseCase{}
	handler := handlers.NewUserHandler(usecase)
	router := setupRouter(handler)

	req, _ := http.NewRequest("GET", "/user/invalid-uuid", nil)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "ID inválido")
}

func TestGetAllUsers_Success(t *testing.T) {
	usecase := &mockUserUseCase{}
	handler := handlers.NewUserHandler(usecase)
	router := setupRouter(handler)

	req, _ := http.NewRequest("GET", "/users", nil)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "nenhum usuário encontrado")
}

func TestGetAllUsers_Empty(t *testing.T) {
	usecase := &mockUserUseCase{}
	handler := handlers.NewUserHandler(usecase)
	router := setupRouter(handler)

	req, _ := http.NewRequest("GET", "/users", nil)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "nenhum usuário encontrado")
}

func TestUpdateUser_Success(t *testing.T) {
	usecase := &mockUserUseCase{}
	handler := handlers.NewUserHandler(usecase)
	router := setupRouter(handler)

	body := map[string]string{
		"id": uuid.New().String(),
	}
	jsonBody, _ := json.Marshal(body)
	userID := uuid.New()
	req, _ := http.NewRequest("PATCH", "/user/update/"+userID.String(), bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), "id", jsonBody))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)	

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "mocked_id")
}

func TestDeleteUser_Success(t *testing.T) {
	usecase := &mockUserUseCase{}
	handler := handlers.NewUserHandler(usecase)
	router := setupRouter(handler)

	body := map[string]string{
		"id": uuid.New().String(),
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("DELETE", "/user/delete", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
 req = req.WithContext(context.WithValue(req.Context(), "id", "mocked_id"))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "user deleted successfully")
}

