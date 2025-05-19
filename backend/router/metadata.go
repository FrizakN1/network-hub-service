package router

import (
	"context"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
)

type Metadata interface {
	setAuthorizationHeader(c *gin.Context) context.Context
}
type DefaultMetadata struct{}

func (m *DefaultMetadata) setAuthorizationHeader(c *gin.Context) context.Context {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(401, gin.H{"error": "authorization header missing"})
		return nil
	}

	md := metadata.New(map[string]string{
		"Authorization": authHeader,
	})

	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

	return ctx
}
