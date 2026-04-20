package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token de autorização necessário"})
			c.Abort()
			return
		}

		secretKey := os.Getenv("JWT_SECRET")
		if secretKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}

		if exp, ok := claims["exp"].(float64); !ok || int64(exp) <= time.Now().Unix() {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expirado"})
			c.Abort()
			return
		}

		email, emailOK := claims["email"].(string)
		userIDRaw, userIDOK := claims["user_id"]
		if !emailOK || email == "" || !userIDOK {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}

		userIDStr := fmt.Sprintf("%v", userIDRaw)
		if _, err := uuid.Parse(userIDStr); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}

		if fp, ok := claims["fp"].(string); ok && fp != "" {
			currentFP := requestFingerprint(c)
			if currentFP != fp {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "sessão inválida para este dispositivo"})
				c.Abort()
				return
			}
		}

		c.Set("user_id", userIDStr)
		c.Set("email", email)

		// Inject plan into context (used by AI and subscription handlers)
		if plan, ok := claims["plan"].(string); ok && plan != "" {
			c.Set("plan", plan)
		} else {
			c.Set("plan", "free")
		}

		c.Next()
	}
}

func requestFingerprint(c *gin.Context) string {
	raw := strings.Join([]string{c.ClientIP(), c.GetHeader("User-Agent"), c.GetHeader("Accept-Language")}, "|")
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
