package router

import (
	"backend/database"
	"backend/utils"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"time"
)

type UserHandler interface {
	handlerEditUser(c *gin.Context)
	handlerCreateUser(c *gin.Context)
	handlerChangeUserStatus(c *gin.Context)
	handlerGetUsers(c *gin.Context)
	handlerGetAuth(c *gin.Context)
	handlerLogout(c *gin.Context)
	handlerLogin(c *gin.Context)
}

type DefaultUserHandler struct {
	UserService database.UserService
}

func (uh *DefaultUserHandler) handlerEditUser(c *gin.Context) {
	sessionHash, ok := c.Get("sessionHash")
	if !ok {
		err := errors.New("сессия не найдена")
		utils.Logger.Println(err)
		handlerError(c, err, 401)
		return
	}

	session := database.GetSession(sessionHash.(string))

	if session.User.Role.Value != "admin" {
		c.JSON(403, nil)
		return
	}

	var (
		err  error
		user database.User
	)

	if err = c.BindJSON(&user); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	if !uh.UserService.ValidateUser(user, "edit") {
		c.JSON(400, nil)
		return
	}

	user.UpdatedAt = sql.NullInt64{
		Int64: time.Now().Unix(),
		Valid: true,
	}

	if err = uh.UserService.EditUser(&user); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	if len(user.Password) != 0 {
		if err = uh.UserService.ChangeUserPassword(&user); err != nil {
			utils.Logger.Println(err)
			handlerError(c, err, 400)
			return
		}

		user.Password = ""
	}

	if err = uh.UserService.DeleteUserSessions(user.ID); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, user)
}

func (uh *DefaultUserHandler) handlerCreateUser(c *gin.Context) {
	sessionHash, ok := c.Get("sessionHash")
	if !ok {
		err := errors.New("сессия не найдена")
		utils.Logger.Println(err)
		handlerError(c, err, 401)
		return
	}

	session := database.GetSession(sessionHash.(string))

	if session.User.Role.Value != "admin" {
		c.JSON(403, nil)
		return
	}

	var (
		err  error
		user database.User
	)

	if err = c.BindJSON(&user); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	if !uh.UserService.ValidateUser(user, "create") {
		c.JSON(400, nil)
		return
	}

	user.CreatedAt = time.Now().Unix()

	if err = uh.UserService.CreateUser(&user); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	user.Password = ""

	c.JSON(200, user)
}

func (uh *DefaultUserHandler) handlerChangeUserStatus(c *gin.Context) {
	sessionHash, ok := c.Get("sessionHash")
	if !ok {
		err := errors.New("сессия не найдена")
		utils.Logger.Println(err)
		handlerError(c, err, 401)
		return
	}

	session := database.GetSession(sessionHash.(string))

	if session.User.Role.Value != "admin" {
		c.JSON(403, nil)
		return
	}

	var user database.User

	if err := c.BindJSON(&user); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	if err := uh.UserService.GetUser(&user); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	if err := uh.UserService.ChangeStatus(&user); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	if err := uh.UserService.DeleteUserSessions(user.ID); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, user)
}

func (uh *DefaultUserHandler) handlerGetUsers(c *gin.Context) {
	sessionHash, ok := c.Get("sessionHash")
	if !ok {
		err := errors.New("сессия не найдена")
		utils.Logger.Println(err)
		handlerError(c, err, 401)
		return
	}

	session := database.GetSession(sessionHash.(string))

	if session.User.Role.Value != "admin" {
		c.JSON(403, nil)
		return
	}

	users, err := uh.UserService.GetUsers()
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, users)
}

func (uh *DefaultUserHandler) handlerGetAuth(c *gin.Context) {
	sessionHash, ok := c.Get("sessionHash")
	if !ok {
		err := errors.New("сессия не найдена")
		utils.Logger.Println(err)
		handlerError(c, err, 401)
		return
	}

	session := database.GetSession(sessionHash.(string))

	c.JSON(200, session.User)
}

func (uh *DefaultUserHandler) handlerLogout(c *gin.Context) {
	sessionHash, ok := c.Get("sessionHash")
	if !ok {
		err := errors.New("сессия не найдена")
		utils.Logger.Println(err)
		handlerError(c, err, 401)
		return
	}

	session := database.GetSession(sessionHash.(string))

	if err := uh.UserService.DeleteUserSessions(session.User.ID); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, nil)
}

func (uh *DefaultUserHandler) handlerLogin(c *gin.Context) {
	var (
		user database.User
		err  error
	)

	if err = c.BindJSON(&user); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	if err = uh.UserService.GetAuthorize(&user); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	if user.ID <= 0 {
		c.JSON(200, gin.H{
			"failure": "Неверный логин/пароль",
		})
		return
	}

	if user.Baned {
		c.JSON(200, gin.H{
			"failure": "Этот аккаунт заблокирован",
		})
		return
	}

	user.Password = ""

	hash, err := uh.UserService.CreateSession(user)
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	token, err := generateToken(hash)
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, gin.H{
		"token": token,
	})
}

func NewUserHandler() UserHandler {
	return &DefaultUserHandler{
		UserService: database.NewUserService(),
	}
}
