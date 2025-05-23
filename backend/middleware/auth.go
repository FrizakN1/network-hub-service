package middleware

import (
	"backend/errors"
	"backend/proto/userpb"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strings"
)

func (m *DefaultMiddleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Error(errors.NewHTTPError(nil, "authorization header is not found", http.StatusUnauthorized))
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Error(errors.NewHTTPError(nil, "invalid token", http.StatusUnauthorized))
			c.Abort()
			return
		}

		tokenString := parts[1]

		ctx := m.Metadata.SetAuthorizationHeader(c)

		res, err := m.UserService.GetSession(ctx, &userpb.GetSessionRequest{Hash: tokenString})
		if err != nil {
			st, _ := status.FromError(err)

			if st.Code() == codes.Unauthenticated {
				c.Error(errors.NewHTTPError(err, "unauthorized", http.StatusUnauthorized))
			} else {
				c.Error(errors.NewHTTPError(err, "failed to get session", http.StatusInternalServerError))
			}

			c.Abort()
			return
		}
		if !res.Exist {
			c.Error(errors.NewHTTPError(err, "unauthorized", http.StatusUnauthorized))
			c.Abort()
			return
		}

		c.Set("session", res.Session)

		c.Next()
	}
}
