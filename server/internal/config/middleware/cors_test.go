package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCORS(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Mock ALLOWED_ORIGINS
	os.Setenv("ALLOWED_ORIGINS", "http://localhost:5173,https://finzen-31s.pages.dev")
	defer os.Unsetenv("ALLOWED_ORIGINS")

	tests := []struct {
		name           string
		origin         string
		expectedOrigin string
	}{
		{
			name:           "Allow localhost",
			origin:         "http://localhost:5173",
			expectedOrigin: "http://localhost:5173",
		},
		{
			name:           "Allow Cloudflare Pages",
			origin:         "https://finzen-31s.pages.dev",
			expectedOrigin: "https://finzen-31s.pages.dev",
		},
		{
			name:           "Deny unknown origin",
			origin:         "https://evil.com",
			expectedOrigin: "",
		},
		{
			name:           "Ignore empty origin",
			origin:         "",
			expectedOrigin: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.Use(CORS())
			r.GET("/test", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			req, _ := http.NewRequest("GET", "/test", nil)
			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
			}
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedOrigin, w.Header().Get("Access-Control-Allow-Origin"))
			if tt.expectedOrigin != "" {
				assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
				assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "PATCH")
			}
		})
	}
}

func TestCORS_Wildcard(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("ALLOWED_ORIGINS", "*")
	defer os.Unsetenv("ALLOWED_ORIGINS")

	r := gin.New()
	r.Use(CORS())
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://any-domain.com")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, "https://any-domain.com", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
}

func TestCORS_Options(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("ALLOWED_ORIGINS", "http://localhost:5173")
	defer os.Unsetenv("ALLOWED_ORIGINS")

	t.Run("Allow options", func(t *testing.T) {
		r := gin.New()
		r.Use(CORS())

		req, _ := http.NewRequest("OPTIONS", "/test", nil)
		req.Header.Set("Origin", "http://localhost:5173")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Equal(t, "http://localhost:5173", w.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("Deny options", func(t *testing.T) {
		r := gin.New()
		r.Use(CORS())

		req, _ := http.NewRequest("OPTIONS", "/test", nil)
		req.Header.Set("Origin", "http://evil.com")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}

