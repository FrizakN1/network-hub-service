package router

import (
	"backend/database"
	"backend/utils"
	"github.com/gin-gonic/gin"
)

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
		"token":    token,
		"userRole": user.Role.Value,
	})
}
