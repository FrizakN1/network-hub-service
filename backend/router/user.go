package router

import (
	"backend/database"
	"backend/utils"
	"errors"
	"github.com/gin-gonic/gin"
)

func handlerGetUsers(c *gin.Context) {
	sessionHash, ok := c.Get("sessionHash")
	if !ok {
		err := errors.New("сессия не найдена")
		utils.Logger.Println(err)
		handlerError(c, err, 400)
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
		handlerError(c, err, 400)
		return
	}

	session := database.GetSession(sessionHash.(string))

	c.JSON(200, session.User)
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
