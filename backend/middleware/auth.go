package middleware

import (
	"backend/errors"
	"backend/proto/userpb"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func (m *DefaultMiddleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Error(errors.NewHTTPError(nil, "authorization header is not found", http.StatusUnauthorized))
			//fmt.Println("Не обнаружен заголовок авторизации")
			//c.JSON(401, gin.H{"error": "Не обнаружен заголовок авторизации"})
			//c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Error(errors.NewHTTPError(nil, "invalid token", http.StatusUnauthorized))
			//fmt.Println("Неверный формат токена")
			//c.JSON(401, gin.H{"error": "Неверный формат токена"})
			//c.Abort()
			return
		}

		tokenString := parts[1]

		ctx := m.Metadata.SetAuthorizationHeader(c)

		res, err := m.UserService.GetSession(ctx, &userpb.GetSessionRequest{Hash: tokenString})
		if err != nil {
			c.Error(errors.NewHTTPError(err, "failed to get session", http.StatusInternalServerError))
			//fmt.Println(err)
			////utils.Logger.Println(err)
			//c.JSON(500, gin.H{"error": "Ошибка при получении сессии"})
			//c.Abort()
			return
		}

		if !res.Exist {
			c.Error(errors.NewHTTPError(err, "unauthorized", http.StatusUnauthorized))
			//fmt.Println("Сессия не найдена")
			//c.JSON(401, gin.H{"error": "Сессия не найдена"})
			//c.Abort()
			return
		}

		c.Set("session", res.Session)

		c.Next()
	}
}
