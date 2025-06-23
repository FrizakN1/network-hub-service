package utils

import (
	"context"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
)

type Metadata interface {
	SetAuthorizationHeader(ctx *gin.Context) context.Context
}
type DefaultMetadata struct{}

func (m *DefaultMetadata) SetAuthorizationHeader(c *gin.Context) context.Context {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(401, gin.H{"error": "authorization header missing"})
		return context.Background()
	}

	md := metadata.New(map[string]string{
		"Authorization": authHeader,
	})

	ctx := metadata.NewOutgoingContext(context.Background(), md)

	return ctx
}
