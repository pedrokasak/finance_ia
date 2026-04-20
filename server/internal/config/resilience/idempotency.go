package resilience

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
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
	StatusCode  int
	Body        json.RawMessage
	ExpiresAt   time.Time
	PayloadHash string
}

var defaultStore = &idempotencyStore{
	store: make(map[string]*cachedResponse),
}

var idempotencyKeyPattern = regexp.MustCompile(`^[a-zA-Z0-9_\-:.]{8,128}$`)

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
		if !idempotencyKeyPattern.MatchString(key) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid Idempotency-Key format"})
			c.Abort()
			return
		}

		body, err := ReadBody(c.Request)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
			c.Abort()
			return
		}
		payloadHash := fmt.Sprintf("%x", sha256.Sum256(body))

		userScope, _ := c.Get("user_id")
		scopedKey := fmt.Sprintf("%v|%s|%s|%s", userScope, c.Request.Method, c.FullPath(), key)

		// Check cache
		defaultStore.mu.RLock()
		cached, found := defaultStore.store[scopedKey]
		defaultStore.mu.RUnlock()

		if found && time.Now().Before(cached.ExpiresAt) {
			if cached.PayloadHash != payloadHash {
				c.JSON(http.StatusConflict, gin.H{"error": "Idempotency-Key reused with different payload"})
				c.Abort()
				return
			}
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
			defaultStore.store[scopedKey] = &cachedResponse{
				StatusCode:  writer.status,
				Body:        json.RawMessage(writer.body.Bytes()),
				ExpiresAt:   time.Now().Add(24 * time.Hour),
				PayloadHash: payloadHash,
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
