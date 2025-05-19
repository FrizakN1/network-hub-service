package router

import (
	"backend/proto/userpb"
	"backend/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log"
	"os"
	"strconv"
)

type UserHandler interface {
	handlerEditUser(c *gin.Context)
	handlerCreateUser(c *gin.Context)
	handlerChangeUserStatus(c *gin.Context)
	handlerGetUsers(c *gin.Context)
	//handlerGetUsersByIds(c *gin.Context, ids []int32) (map[int32]*userpb.User, error)
}

type DefaultUserHandler struct {
	Metadata    Metadata
	Privilege   Privilege
	UserService userpb.UserServiceClient
}

func NewUserHandler(userService *userpb.UserServiceClient) UserHandler {
	return &DefaultUserHandler{
		UserService: *userService,
		Privilege:   &DefaultPrivilege{},
		Metadata:    &DefaultMetadata{},
	}
}

func InitUserClient() userpb.UserServiceClient {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", os.Getenv("USER_SERVICE_ADDRESS"), os.Getenv("USER_SERVICE_PORT")), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect to user service: %v", err)
	}

	userClient := userpb.NewUserServiceClient(conn)

	return userClient
}

//func (h DefaultUserHandler) handlerGetUsersByIds(c *gin.Context, ids []int32) (map[int32]*userpb.User, error) {
//	ctx := h.Metadata.setAuthorizationHeader(c)
//
//	userResp, err := h.UserService.GetUsersByIds(ctx, &userpb.GetUsersByIdsRequest{Ids: ids})
//	if err != nil {
//		utils.Logger.Println(err)
//		handlerError(c, err, 500)
//		return nil, err
//	}
//
//	usersMap := make(map[int32]*userpb.User)
//	for _, user := range userResp.Users {
//		usersMap[user.Id] = user
//	}
//
//	return usersMap, nil
//}

func (h DefaultUserHandler) handlerGetUsers(c *gin.Context) {
	_, _, isOperatorOrHigher := h.Privilege.getPrivilege(c)

	if !isOperatorOrHigher {
		c.JSON(403, nil)
		return
	}

	ctx := h.Metadata.setAuthorizationHeader(c)

	res, err := h.UserService.GetUsers(ctx, &userpb.Empty{})
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 500)
		return
	}

	c.JSON(200, res.Users)
}

func (h DefaultUserHandler) handlerCreateUser(c *gin.Context) {
	_, _, isOperatorOrHigher := h.Privilege.getPrivilege(c)

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

	ctx := h.Metadata.setAuthorizationHeader(c)

	res, err := h.UserService.CreateUser(ctx, user)
	if err != nil {
		utils.Logger.Println(err)
		c.JSON(500, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(200, res)
}

func (h DefaultUserHandler) handlerEditUser(c *gin.Context) {
	_, _, isOperatorOrHigher := h.Privilege.getPrivilege(c)

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

	ctx := h.Metadata.setAuthorizationHeader(c)

	res, err := h.UserService.EditUser(ctx, user)
	if err != nil {
		utils.Logger.Println(err)
		c.JSON(500, gin.H{"error": "failed to edit user"})
		return
	}

	c.JSON(200, res)
}

func (h DefaultUserHandler) handlerChangeUserStatus(c *gin.Context) {
	_, _, isOperatorOrHigher := h.Privilege.getPrivilege(c)

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

	ctx := h.Metadata.setAuthorizationHeader(c)

	res, err := h.UserService.ChangeUserStatus(ctx, &userpb.ChangeUserStatusRequest{Id: int32(userID)})
	if err != nil {
		utils.Logger.Println(err)
		c.JSON(500, gin.H{"error": "failed to change user status"})
		return
	}

	c.JSON(200, res)
}
