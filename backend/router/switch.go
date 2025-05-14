package router

import (
	"backend/database"
	"backend/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"time"
)

type SwitchHandler interface {
	handlerGetSwitches(c *gin.Context)
	handlerEditSwitch(c *gin.Context)
	handlerCreateSwitch(c *gin.Context)
}

type DefaultSwitchHandler struct {
	SwitchService database.SwitchService
}

func (sh *DefaultSwitchHandler) handlerGetSwitches(c *gin.Context) {
	switches, err := sh.SwitchService.GetSwitches()
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, switches)
}

func (sh *DefaultSwitchHandler) handlerEditSwitch(c *gin.Context) {
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

	var _switch database.Switch

	if err := c.BindJSON(&_switch); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	if len(_switch.Name) == 0 || _switch.PortAmount == 0 {
		c.JSON(400, nil)
		return
	}

	if err := sh.SwitchService.EditSwitch(&_switch); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, _switch)
}

func (sh *DefaultSwitchHandler) handlerCreateSwitch(c *gin.Context) {
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

	var _switch database.Switch

	if err := c.BindJSON(&_switch); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	if len(_switch.Name) == 0 || _switch.PortAmount == 0 {
		c.JSON(400, nil)
		return
	}

	_switch.CreatedAt = time.Now().Unix()

	if err := sh.SwitchService.CreateSwitch(&_switch); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, _switch)
}

func NewSwitchHandler() SwitchHandler {
	return &DefaultSwitchHandler{
		SwitchService: &database.DefaultSwitchService{},
	}
}
