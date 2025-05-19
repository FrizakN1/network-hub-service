package router

import (
	"backend/proto/userpb"
	"backend/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

type AuthHandler interface {
	authMiddleware() gin.HandlerFunc
	handlerGetAuth(c *gin.Context)
	handlerLogout(c *gin.Context)
	handlerLogin(c *gin.Context)
}

type DefaultAuthHandler struct {
	Metadata    Metadata
	Privilege   Privilege
	UserService userpb.UserServiceClient
}

func NewAuthHandler(userClient *userpb.UserServiceClient) AuthHandler {
	return &DefaultAuthHandler{
		Metadata:    &DefaultMetadata{},
		Privilege:   &DefaultPrivilege{},
		UserService: *userClient,
	}
}

func (h *DefaultAuthHandler) handlerGetAuth(c *gin.Context) {
	session, _, _ := h.Privilege.getPrivilege(c)

	c.JSON(200, session.User)
}

func (h *DefaultAuthHandler) handlerLogout(c *gin.Context) {
	session, _, _ := h.Privilege.getPrivilege(c)

	ctx := h.Metadata.setAuthorizationHeader(c)

	_, err := h.UserService.Logout(ctx, &userpb.LogoutRequest{Hash: session.Hash})
	if err != nil {
		utils.Logger.Println(err)
		c.JSON(500, gin.H{"error": "failed to logout"})
		return
	}

	c.JSON(200, nil)
}

func (h *DefaultAuthHandler) handlerLogin(c *gin.Context) {
	var err error
	loginData := &userpb.LoginRequest{}

	if err = c.BindJSON(&loginData); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	ctx := h.Metadata.setAuthorizationHeader(c)

	res, err := h.UserService.Login(ctx, loginData)
	if err != nil {
		utils.Logger.Println(err)
		c.JSON(500, gin.H{"error": "failed to login"})
		return
	}

	c.JSON(200, gin.H{
		"token":   res.Hash,
		"failure": res.Failure,
	})
}

func (h *DefaultAuthHandler) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			fmt.Println("Не обнаружен заголовок авторизации")
			c.JSON(401, gin.H{"error": "Не обнаружен заголовок авторизации"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			fmt.Println("Неверный формат токена")
			c.JSON(401, gin.H{"error": "Неверный формат токена"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		ctx := h.Metadata.setAuthorizationHeader(c)

		res, err := h.UserService.GetSession(ctx, &userpb.GetSessionRequest{Hash: tokenString})
		if err != nil {
			fmt.Println(err)
			utils.Logger.Println(err)
			c.JSON(500, gin.H{"error": "Ошибка при получении сессии"})
			c.Abort()
			return
		}

		if !res.Exist {
			fmt.Println("Сессия не найдена")
			c.JSON(401, gin.H{"error": "Сессия не найдена"})
			c.Abort()
			return
		}

		c.Set("session", res.Session)

		c.Next()
	}
}
