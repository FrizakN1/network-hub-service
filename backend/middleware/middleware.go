package middleware

import (
	"backend/proto/userpb"
	"backend/utils"
	"github.com/gin-gonic/gin"
)

type Middleware interface {
	AuthMiddleware() gin.HandlerFunc
	CorsMiddleware() gin.HandlerFunc
	ErrorMiddleware() gin.HandlerFunc
}

type DefaultMiddleware struct {
	Metadata    utils.Metadata
	UserService userpb.UserServiceClient
	Logger      utils.Logger
}

func NewMiddleware(userClient *userpb.UserServiceClient, logger *utils.Logger) Middleware {
	return &DefaultMiddleware{
		Metadata:    &utils.DefaultMetadata{},
		UserService: *userClient,
		Logger:      *logger,
	}
}
