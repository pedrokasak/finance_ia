package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetUserID extrai o ID do usuário do contexto do Gin de forma segura.
// Retorna uuid.Nil caso o ID não esteja presente ou seja inválido.
func GetUserID(c *gin.Context) uuid.UUID {
	val, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil
	}

	// Tenta converter de string (comum se vier do JWT)
	if s, ok := val.(string); ok {
		id, err := uuid.Parse(s)
		if err != nil {
			return uuid.Nil
		}
		return id
	}

	// Tenta converter diretamente de uuid.UUID
	if id, ok := val.(uuid.UUID); ok {
		return id
	}

	return uuid.Nil
}
