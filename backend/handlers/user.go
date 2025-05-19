package handlers

import (
	"backend/errors"
	"backend/proto/userpb"
	"backend/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"os"
	"strconv"
)

type UserHandler interface {
	HandlerEditUser(c *gin.Context)
	HandlerCreateUser(c *gin.Context)
	HandlerChangeUserStatus(c *gin.Context)
	HandlerGetUsers(c *gin.Context)
}

type DefaultUserHandler struct {
	Metadata    utils.Metadata
	Privilege   Privilege
	UserService userpb.UserServiceClient
	Logger      utils.Logger
}

func NewUserHandler(userService *userpb.UserServiceClient) UserHandler {
	return &DefaultUserHandler{
		UserService: *userService,
		Privilege:   &DefaultPrivilege{},
		Metadata:    &utils.DefaultMetadata{},
	}
}

func InitUserClient() userpb.UserServiceClient {
	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%s", os.Getenv("USER_SERVICE_ADDRESS"), os.Getenv("USER_SERVICE_PORT")),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("could not connect to user service: %v", err)
	}

	userClient := userpb.NewUserServiceClient(conn)

	return userClient
}

func (h DefaultUserHandler) HandlerGetUsers(c *gin.Context) {
	_, _, isOperatorOrHigher := h.Privilege.getPrivilege(c)

	if !isOperatorOrHigher {
		c.Error(errors.NewHTTPError(nil, "forbidden", http.StatusForbidden))
		return
	}

	ctx := h.Metadata.SetAuthorizationHeader(c)

	res, err := h.UserService.GetUsers(ctx, &userpb.Empty{})
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to get users", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, res.Users)
}

func (h DefaultUserHandler) HandlerCreateUser(c *gin.Context) {
	_, _, isOperatorOrHigher := h.Privilege.getPrivilege(c)

	if !isOperatorOrHigher {
		c.Error(errors.NewHTTPError(nil, "forbidden", http.StatusForbidden))
		return
	}

	user := &userpb.CreateUserRequest{}

	if err := c.BindJSON(&user); err != nil {
		c.Error(errors.NewHTTPError(err, "invalid json", http.StatusBadRequest))
		return
	}

	ctx := h.Metadata.SetAuthorizationHeader(c)

	res, err := h.UserService.CreateUser(ctx, user)
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to create user", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h DefaultUserHandler) HandlerEditUser(c *gin.Context) {
	_, _, isOperatorOrHigher := h.Privilege.getPrivilege(c)

	if !isOperatorOrHigher {
		c.Error(errors.NewHTTPError(nil, "forbidden", http.StatusForbidden))
		return
	}

	user := &userpb.EditUserRequest{}

	if err := c.BindJSON(&user); err != nil {
		c.Error(errors.NewHTTPError(err, "invalid json", http.StatusBadRequest))
		return
	}

	ctx := h.Metadata.SetAuthorizationHeader(c)

	res, err := h.UserService.EditUser(ctx, user)
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to edit user", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h DefaultUserHandler) HandlerChangeUserStatus(c *gin.Context) {
	_, _, isOperatorOrHigher := h.Privilege.getPrivilege(c)

	if !isOperatorOrHigher {
		c.Error(errors.NewHTTPError(nil, "forbidden", http.StatusForbidden))
		return
	}

	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to parse param(id) to int", http.StatusBadRequest))
		return
	}

	ctx := h.Metadata.SetAuthorizationHeader(c)

	res, err := h.UserService.ChangeUserStatus(ctx, &userpb.ChangeUserStatusRequest{Id: int32(userID)})
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to change user status", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, res)
}
