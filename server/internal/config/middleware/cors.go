package middleware

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		allowedOriginsRaw := os.Getenv("ALLOWED_ORIGINS")
		origin := c.Request.Header.Get("Origin")

		if origin == "" {
			c.Next()
			return
		}

		isAllowed := false
		if allowedOriginsRaw == "*" || allowedOriginsRaw == "" {
			isAllowed = true
		} else {
			allowedOrigins := strings.Split(allowedOriginsRaw, ",")
			for _, o := range allowedOrigins {
				trimmed := strings.TrimSpace(o)
				if trimmed == "*" || trimmed == origin {
					isAllowed = true
					break
				}
			}
		}

		if isAllowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, Idempotency-Key")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH, HEAD")
		}

		if c.Request.Method == "OPTIONS" {
			if isAllowed {
				c.AbortWithStatus(204)
			} else {
				c.AbortWithStatus(403)
			}
			return
		}

		c.Next()
	}
}
