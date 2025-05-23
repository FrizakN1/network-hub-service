package middleware

import (
	httpErrors "backend/errors"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (m *DefaultMiddleware) ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				m.Logger.Println(err)

				var httpErr *httpErrors.HTTPError
				if errors.As(err.Err, &httpErr) {
					c.AbortWithStatusJSON(httpErr.Code, gin.H{"error": httpErr.Message})
					return
				}

				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
				return
			}
		}
	}
}
