package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func (m *DefaultMiddleware) CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Проверка ссылок на соответствие
		if os.Getenv("ALLOW_ORIGIN") == origin {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(http.StatusNoContent)
				return
			} else {
				c.Next()
			}
		}
	}
}
