package auth_test

import (
	"bytes"
	"encoding/json"
	"errors"
	domainAuth "finance-ia/internal/domain/auth"
	handlersAuth "finance-ia/internal/interfaces/handlers/auth"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Mock da camada de usecase
type mockUserUseCase struct{}

func (m *mockUserUseCase) Login(email, password, code string) (string, error) {
	if email == "invalid@example.com" {
		return "", errors.New("invalid credentials")
	}
	if email == "2fa@example.com" && code == "" {
		return "", errors.New("2fa_required")
	}
	return "mocked_token", nil
}

func (m *mockUserUseCase) Setup2FA(email string) (string, string, error) {
	return "mock_secret", "mock_url", nil
}

func (m *mockUserUseCase) Verify2FA(email, code string) error {
	if code == "invalid" {
		return errors.New("invalid code")
	}
	return nil
}

func (m *mockUserUseCase) Disable2FA(email string) error {
	return nil
}

func (m *mockUserUseCase) Update(u *domainAuth.Authentication) error     { return nil }
func (m *mockUserUseCase) ForgotPassword(email string) error             { return nil }
func (m *mockUserUseCase) ResetPassword(token, newPassword string) error { return nil }
func (m *mockUserUseCase) GetByEmail(email string) (*domainAuth.Authentication, error) {
	return &domainAuth.Authentication{Email: email}, nil
}
func (m *mockUserUseCase) Delete(u *domainAuth.Authentication) error { return nil }
func (m *mockUserUseCase) Logout(token string) error                 { return nil }
func setupRouter(handler *handlersAuth.AuthHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/auth/login", handler.Login)
	r.POST("/auth/forgot-password", handler.ForgotPassword)
	r.POST("/auth/reset-password", handler.ResetPassword)
	r.POST("/auth/logout", handler.Logout)
	return r
}

func TestLogin_Success(t *testing.T) {
	usecase := &mockUserUseCase{}
	handler := handlersAuth.NewAuthHandler(usecase)
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
	handler := handlersAuth.NewAuthHandler(usecase)
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

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	assert.Contains(t, resp.Body.String(), "invalid credentials")
}

func TestLogin_InvalidRequestBody(t *testing.T) {
	usecase := &mockUserUseCase{}
	handler := handlersAuth.NewAuthHandler(usecase)
	router := setupRouter(handler)

	body := "invalid json"
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "invalid character")
}

func TestRecoveryPassword(t *testing.T) {
	usecase := &mockUserUseCase{}
	handler := handlersAuth.NewAuthHandler(usecase)
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
	assert.Contains(t, resp.Body.String(), "Se o email existir")
}

func TestResetPassword(t *testing.T) {
	usecase := &mockUserUseCase{}
	handler := handlersAuth.NewAuthHandler(usecase)
	router := setupRouter(handler)

	body := map[string]string{
		"token":        "valid_token",
		"new_password": "newpassword123",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/auth/reset-password", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "Senha redefinida com sucesso!")
}

func TestResetPassword_InvalidRequestBody(t *testing.T) {
	usecase := &mockUserUseCase{}
	handler := handlersAuth.NewAuthHandler(usecase)
	router := setupRouter(handler)

	body := map[string]string{
		"token":        "invalid_token",
		"new_password": "wrong",
	}
	req, _ := http.NewRequest("POST", "/auth/reset-password", bytes.NewBufferString(body["token"]))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "Dados inválidos")
}
