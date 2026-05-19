package resilience

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func newIdempotencyTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", c.GetHeader("X-Test-User-ID"))
		c.Next()
	})
	r.Use(IdempotencyMiddleware())
	r.POST("/resource", func(c *gin.Context) {
		var body map[string]any
		_ = c.ShouldBindJSON(&body)
		c.JSON(http.StatusCreated, gin.H{"ok": true, "body": body})
	})
	return r
}

func TestIdempotency_SameKeySamePayload_ReturnsCached(t *testing.T) {
	r := newIdempotencyTestRouter()
	payload := []byte(`{"value":1}`)

	req1, _ := http.NewRequest("POST", "/resource", bytes.NewBuffer(payload))
	req1.Header.Set("Content-Type", "application/json")
	req1.Header.Set("Idempotency-Key", "abc12345-test-key")
	req1.Header.Set("X-Test-User-ID", "user-a")
	res1 := httptest.NewRecorder()
	r.ServeHTTP(res1, req1)
	assert.Equal(t, http.StatusCreated, res1.Code)

	req2, _ := http.NewRequest("POST", "/resource", bytes.NewBuffer(payload))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("Idempotency-Key", "abc12345-test-key")
	req2.Header.Set("X-Test-User-ID", "user-a")
	res2 := httptest.NewRecorder()
	r.ServeHTTP(res2, req2)
	assert.Equal(t, http.StatusCreated, res2.Code)
	assert.JSONEq(t, res1.Body.String(), res2.Body.String())
}

func TestIdempotency_SameKeyDifferentPayload_ReturnsConflict(t *testing.T) {
	r := newIdempotencyTestRouter()

	req1, _ := http.NewRequest("POST", "/resource", bytes.NewBuffer([]byte(`{"value":1}`)))
	req1.Header.Set("Content-Type", "application/json")
	req1.Header.Set("Idempotency-Key", "abc12345-test-key-2")
	req1.Header.Set("X-Test-User-ID", "user-a")
	res1 := httptest.NewRecorder()
	r.ServeHTTP(res1, req1)
	assert.Equal(t, http.StatusCreated, res1.Code)

	req2, _ := http.NewRequest("POST", "/resource", bytes.NewBuffer([]byte(`{"value":2}`)))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("Idempotency-Key", "abc12345-test-key-2")
	req2.Header.Set("X-Test-User-ID", "user-a")
	res2 := httptest.NewRecorder()
	r.ServeHTTP(res2, req2)
	assert.Equal(t, http.StatusConflict, res2.Code)

	var body map[string]string
	_ = json.Unmarshal(res2.Body.Bytes(), &body)
	assert.Contains(t, body["error"], "Idempotency-Key")
}

func TestIdempotency_SameKeyDifferentUser_IsolatedScope(t *testing.T) {
	r := newIdempotencyTestRouter()

	req1, _ := http.NewRequest("POST", "/resource", bytes.NewBuffer([]byte(`{"value":1}`)))
	req1.Header.Set("Content-Type", "application/json")
	req1.Header.Set("Idempotency-Key", "abc12345-test-key-3")
	req1.Header.Set("X-Test-User-ID", "user-a")
	res1 := httptest.NewRecorder()
	r.ServeHTTP(res1, req1)
	assert.Equal(t, http.StatusCreated, res1.Code)

	req2, _ := http.NewRequest("POST", "/resource", bytes.NewBuffer([]byte(`{"value":2}`)))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("Idempotency-Key", "abc12345-test-key-3")
	req2.Header.Set("X-Test-User-ID", "user-b")
	res2 := httptest.NewRecorder()
	r.ServeHTTP(res2, req2)
	assert.Equal(t, http.StatusCreated, res2.Code)
	assert.NotEqual(t, res1.Body.String(), res2.Body.String())
}
