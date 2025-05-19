package handlers

import (
	"backend/errors"
	"backend/proto/userpb"
	"backend/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthHandler interface {
	HandlerGetAuth(c *gin.Context)
	HandlerLogout(c *gin.Context)
	HandlerLogin(c *gin.Context)
}

type DefaultAuthHandler struct {
	Metadata    utils.Metadata
	Privilege   Privilege
	UserService userpb.UserServiceClient
}

func NewAuthHandler(userClient *userpb.UserServiceClient) AuthHandler {
	return &DefaultAuthHandler{
		Metadata:    &utils.DefaultMetadata{},
		Privilege:   &DefaultPrivilege{},
		UserService: *userClient,
	}
}

func (h *DefaultAuthHandler) HandlerGetAuth(c *gin.Context) {
	session, _, _ := h.Privilege.getPrivilege(c)

	c.JSON(http.StatusOK, session.User)
}

func (h *DefaultAuthHandler) HandlerLogout(c *gin.Context) {
	session, _, _ := h.Privilege.getPrivilege(c)

	ctx := h.Metadata.SetAuthorizationHeader(c)

	_, err := h.UserService.Logout(ctx, &userpb.LogoutRequest{Hash: session.Hash})
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to logout", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (h *DefaultAuthHandler) HandlerLogin(c *gin.Context) {
	var err error
	loginData := &userpb.LoginRequest{}

	if err = c.BindJSON(&loginData); err != nil {
		c.Error(errors.NewHTTPError(err, "invalid json", http.StatusBadRequest))
		return
	}

	ctx := h.Metadata.SetAuthorizationHeader(c)

	res, err := h.UserService.Login(ctx, loginData)
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to login", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":   res.Hash,
		"failure": res.Failure,
	})
}
