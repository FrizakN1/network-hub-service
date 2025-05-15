package router

import (
	"backend/database"
	"backend/utils"
	"github.com/gin-gonic/gin"
	"time"
)

func (h *DefaultHandler) handlerGetSwitches(c *gin.Context) {
	switches, err := h.SwitchService.GetSwitches()
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, switches)
}

func (h *DefaultHandler) handlerEditSwitch(c *gin.Context) {
	_, _, isOperatorOrHigher := h.getPrivilege(c)

	if !isOperatorOrHigher {
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

	if err := h.SwitchService.EditSwitch(&_switch); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, _switch)
}

func (h *DefaultHandler) handlerCreateSwitch(c *gin.Context) {
	_, _, isOperatorOrHigher := h.getPrivilege(c)

	if !isOperatorOrHigher {
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

	if err := h.SwitchService.CreateSwitch(&_switch); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, _switch)
}
