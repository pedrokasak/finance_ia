package resilience

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// idempotencyStore is a simple in-memory store for idempotency keys
// In production, replace with Redis for distributed support
type idempotencyStore struct {
	mu    sync.RWMutex
	store map[string]*cachedResponse
}

type cachedResponse struct {
	StatusCode int
	Body       json.RawMessage
	ExpiresAt  time.Time
}

var defaultStore = &idempotencyStore{
	store: make(map[string]*cachedResponse),
}

// IdempotencyMiddleware validates and caches responses by Idempotency-Key header.
// If a request with the same key was already processed, the cached response is returned.
func IdempotencyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only apply to mutating methods
		if c.Request.Method == http.MethodGet || c.Request.Method == http.MethodDelete {
			c.Next()
			return
		}

		key := c.GetHeader("Idempotency-Key")
		if key == "" {
			c.Next()
			return
		}

		// Check cache
		defaultStore.mu.RLock()
		cached, found := defaultStore.store[key]
		defaultStore.mu.RUnlock()

		if found && time.Now().Before(cached.ExpiresAt) {
			c.Data(cached.StatusCode, "application/json", cached.Body)
			c.Abort()
			return
		}

		// Wrap response writer to capture response
		writer := &responseCapture{ResponseWriter: c.Writer, body: &bytes.Buffer{}}
		c.Writer = writer

		c.Next()

		// Cache the response for 24 hours
		if writer.status >= 200 && writer.status < 300 {
			defaultStore.mu.Lock()
			defaultStore.store[key] = &cachedResponse{
				StatusCode: writer.status,
				Body:       json.RawMessage(writer.body.Bytes()),
				ExpiresAt:  time.Now().Add(24 * time.Hour),
			}
			defaultStore.mu.Unlock()
		}

		// Periodic cleanup of expired keys
		go cleanupExpiredKeys()
	}
}

type responseCapture struct {
	gin.ResponseWriter
	body   *bytes.Buffer
	status int
}

func (r *responseCapture) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func (r *responseCapture) WriteHeader(statusCode int) {
	r.status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *responseCapture) WriteString(s string) (int, error) {
	return r.Write([]byte(s))
}

// ReadBody reads and restores the request body (useful for logging)
func ReadBody(r *http.Request) ([]byte, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	r.Body = io.NopCloser(bytes.NewBuffer(body))
	return body, nil
}

func cleanupExpiredKeys() {
	defaultStore.mu.Lock()
	defer defaultStore.mu.Unlock()
	now := time.Now()
	for k, v := range defaultStore.store {
		if now.After(v.ExpiresAt) {
			delete(defaultStore.store, k)
		}
	}
}
