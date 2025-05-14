package router

import (
	"backend/proto/userpb"
	"backend/utils"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log"
	"strconv"
)

var userClient userpb.UserServiceClient

func InitUserClient() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect to user service: %v", err)
	}

	userClient = userpb.NewUserServiceClient(conn)
}

func (h *DefaultHandler) handlerGetUsers(c *gin.Context) {
	res, err := userClient.GetUsers(context.Background(), &userpb.Empty{})
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to get users"})
		return
	}

	c.JSON(200, res.Users)
}

func (h *DefaultHandler) handlerCreateUser(c *gin.Context) {
	session, ok := c.Get("session")
	if !ok {
		err := errors.New("сессия не найдена")
		utils.Logger.Println(err)
		handlerError(c, err, 401)
		return
	}

	if session.(userpb.Session).User.Role.Value != "admin" && session.(userpb.Session).User.Role.Value != "operator" {
		c.JSON(403, nil)
		return
	}

	user := &userpb.CreateUserRequest{}

	if err := c.BindJSON(&user); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	res, err := userClient.CreateUser(context.Background(), user)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(200, res)
}

func (h *DefaultHandler) handlerEditUser(c *gin.Context) {
	session, ok := c.Get("session")
	if !ok {
		err := errors.New("сессия не найдена")
		utils.Logger.Println(err)
		handlerError(c, err, 401)
		return
	}

	if session.(userpb.Session).User.Role.Value != "admin" && session.(userpb.Session).User.Role.Value != "operator" {
		c.JSON(403, nil)
		return
	}

	user := &userpb.EditUserRequest{}

	if err := c.BindJSON(&user); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	res, err := userClient.EditUser(context.Background(), user)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to edit user"})
		return
	}

	c.JSON(200, res)
}

func (h *DefaultHandler) handlerChangeUserStatus(c *gin.Context) {
	session, ok := c.Get("session")
	if !ok {
		err := errors.New("сессия не найдена")
		utils.Logger.Println(err)
		handlerError(c, err, 401)
		return
	}

	if session.(userpb.Session).User.Role.Value != "admin" && session.(userpb.Session).User.Role.Value != "operator" {
		c.JSON(403, nil)
		return
	}

	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	res, err := userClient.ChangeUserStatus(context.Background(), &userpb.ChangeUserStatusRequest{Id: int32(userID)})
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to change user status"})
		return
	}

	c.JSON(200, res)
}

func (h *DefaultHandler) handlerGetAuth(c *gin.Context) {
	session, ok := c.Get("session")
	if !ok {
		err := errors.New("сессия не найдена")
		utils.Logger.Println(err)
		handlerError(c, err, 401)
		return
	}

	c.JSON(200, session.(*userpb.Session).User)
}

func (h *DefaultHandler) handlerLogout(c *gin.Context) {
	session, ok := c.Get("session")
	if !ok {
		err := errors.New("сессия не найдена")
		utils.Logger.Println(err)
		handlerError(c, err, 401)
		return
	}

	_, err := userClient.Logout(context.Background(), &userpb.LogoutRequest{Hash: session.(*userpb.Session).Hash})
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to logout"})
		return
	}

	c.JSON(200, nil)
}

func (uh *DefaultHandler) handlerLogin(c *gin.Context) {
	var (
		loginData userpb.LoginRequest
		err       error
	)

	if err = c.BindJSON(&loginData); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	res, err := userClient.Login(context.Background(), &loginData)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to login"})
		return
	}

	if res.Failure != "" {
		c.JSON(200, gin.H{
			"failure": res.Failure,
		})

		return
	}

	c.JSON(200, gin.H{
		"token": res.Hash,
	})
}
