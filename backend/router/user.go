package router

import (
	"backend/proto/userpb"
	"backend/utils"
	"fmt"
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

func (h *DefaultHandler) handlerGetUsersByIds(c *gin.Context, ids []int32) (map[int32]*userpb.User, error) {
	userResp, err := userClient.GetUsersByIds(c, &userpb.GetUsersByIdsRequest{Ids: ids})
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 500)
		return nil, err
	}

	usersMap := make(map[int32]*userpb.User)
	for _, user := range userResp.Users {
		usersMap[user.Id] = user
	}

	return usersMap, nil
}

func (h *DefaultHandler) handlerGetUsers(c *gin.Context) {
	_, _, isOperatorOrHigher := h.getPrivilege(c)

	if !isOperatorOrHigher {
		c.JSON(403, nil)
		return
	}

	res, err := userClient.GetUsers(c.Request.Context(), &userpb.Empty{})
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 500)
		return
	}

	c.JSON(200, res.Users)
}

func (h *DefaultHandler) handlerCreateUser(c *gin.Context) {
	_, _, isOperatorOrHigher := h.getPrivilege(c)

	if !isOperatorOrHigher {
		c.JSON(403, nil)
		return
	}

	user := &userpb.CreateUserRequest{}

	if err := c.BindJSON(&user); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	res, err := userClient.CreateUser(c.Request.Context(), user)
	if err != nil {
		utils.Logger.Println(err)
		c.JSON(500, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(200, res)
}

func (h *DefaultHandler) handlerEditUser(c *gin.Context) {
	_, _, isOperatorOrHigher := h.getPrivilege(c)

	if !isOperatorOrHigher {
		c.JSON(403, nil)
		return
	}

	user := &userpb.EditUserRequest{}

	if err := c.BindJSON(&user); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	res, err := userClient.EditUser(c.Request.Context(), user)
	if err != nil {
		utils.Logger.Println(err)
		c.JSON(500, gin.H{"error": "failed to edit user"})
		return
	}

	c.JSON(200, res)
}

func (h *DefaultHandler) handlerChangeUserStatus(c *gin.Context) {
	_, _, isOperatorOrHigher := h.getPrivilege(c)

	if !isOperatorOrHigher {
		c.JSON(403, nil)
		return
	}

	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	res, err := userClient.ChangeUserStatus(c.Request.Context(), &userpb.ChangeUserStatusRequest{Id: int32(userID)})
	if err != nil {
		utils.Logger.Println(err)
		c.JSON(500, gin.H{"error": "failed to change user status"})
		return
	}

	c.JSON(200, res)
}

func (h *DefaultHandler) handlerGetAuth(c *gin.Context) {
	fmt.Println(1233213123)
	session, _, _ := h.getPrivilege(c)

	c.JSON(200, session.User)
}

func (h *DefaultHandler) handlerLogout(c *gin.Context) {
	session, _, _ := h.getPrivilege(c)

	_, err := userClient.Logout(c.Request.Context(), &userpb.LogoutRequest{Hash: session.Hash})
	if err != nil {
		utils.Logger.Println(err)
		c.JSON(500, gin.H{"error": "failed to logout"})
		return
	}

	c.JSON(200, nil)
}

func (h *DefaultHandler) handlerLogin(c *gin.Context) {
	var err error
	loginData := &userpb.LoginRequest{}

	if err = c.BindJSON(&loginData); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	res, err := userClient.Login(c.Request.Context(), loginData)
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
