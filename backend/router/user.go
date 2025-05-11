package router

import (
	"backend/database"
	"backend/utils"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"time"
)

func handlerEditUser(c *gin.Context) {
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

	if !user.ValidateUser("edit") {
		c.JSON(400, nil)
		return
	}

	user.UpdatedAt = sql.NullInt64{
		Int64: time.Now().Unix(),
		Valid: true,
	}

	if err = user.EditUser(); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	if len(user.Password) != 0 {
		user.Password, err = utils.Encrypt(user.Password)
		if err != nil {
			utils.Logger.Println(err)
			handlerError(c, err, 400)
			return
		}

		if err = user.ChangeUserPassword(); err != nil {
			utils.Logger.Println(err)
			handlerError(c, err, 400)
			return
		}

		user.Password = ""
	}

	if err = database.DeleteUserSessions(user.ID); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, user)
}

func handlerCreateUser(c *gin.Context) {
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

	if !user.ValidateUser("create") {
		c.JSON(400, nil)
		return
	}

	user.CreatedAt = time.Now().Unix()
	user.Password, err = utils.Encrypt(user.Password)
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	if err = user.CreateUser(); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	user.Password = ""

	c.JSON(200, user)
}

func handlerChangeUserStatus(c *gin.Context) {
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

	if err := user.GetUser(); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	if err := user.ChangeStatus(); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	if err := database.DeleteUserSessions(user.ID); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, user)
}

func handlerGetUsers(c *gin.Context) {
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

	users, err := database.GetUsers()
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, users)
}

func handlerGetAuth(c *gin.Context) {
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

func handlerLogout(c *gin.Context) {
	sessionHash, ok := c.Get("sessionHash")
	if !ok {
		err := errors.New("сессия не найдена")
		utils.Logger.Println(err)
		handlerError(c, err, 401)
		return
	}

	session := database.GetSession(sessionHash.(string))

	if err := database.DeleteUserSessions(session.User.ID); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, nil)
}

func handlerLogin(c *gin.Context) {
	var (
		user database.User
		err  error
	)

	if err = c.BindJSON(&user); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	user.Password, err = utils.Encrypt(user.Password)
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	if err = user.GetAuthorize(); err != nil {
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

	hash, err := database.CreateSession(user)
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
