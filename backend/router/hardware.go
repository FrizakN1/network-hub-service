package router

import (
	"backend/database"
	"backend/utils"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

func handlerGetHardwareFiles(c *gin.Context) {
	hardwareID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	files, err := database.GetHardwareFiles(hardwareID)
	if err != nil {
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, files)
}

func handlerGetHardwareByID(c *gin.Context) {
	var (
		err      error
		hardware database.Hardware
	)

	hardware.ID, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	if err = hardware.GetHardwareByID(); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, hardware)
}

func handlerEditHardware(c *gin.Context) {
	var hardware database.Hardware

	if err := c.BindJSON(&hardware); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	if !hardware.ValidateHardware() {
		c.JSON(400, nil)
		return
	}

	hardware.UpdatedAt = sql.NullInt64{Int64: time.Now().Unix(), Valid: true}

	if err := hardware.EditHardware(); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, hardware)
}

func handlerCreateHardware(c *gin.Context) {
	var hardware database.Hardware

	if err := c.BindJSON(&hardware); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	if !hardware.ValidateHardware() {
		c.JSON(400, nil)
		return
	}

	hardware.CreatedAt = time.Now().Unix()

	if err := hardware.CreateHardware(); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, hardware)
}

func handlerGetSearchHardware(c *gin.Context) {
	request := struct {
		Text   string
		Offset int
	}{}

	if err := c.BindJSON(&request); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	hardware, count, err := database.GetSearchHardware(request.Text, request.Offset)
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, gin.H{
		"Hardware": hardware,
		"Count":    count,
	})
}

func handlerGetNodeHardware(c *gin.Context) {
	request := struct {
		Offset int
	}{}

	if err := c.BindJSON(&request); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	nodeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	hardware, count, err := database.GetNodeHardware(nodeID, request.Offset)
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, gin.H{
		"Hardware": hardware,
		"Count":    count,
	})
}

func handlerGetHouseHardware(c *gin.Context) {
	request := struct {
		Offset int
	}{}

	if err := c.BindJSON(&request); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	houseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	hardware, count, err := database.GetHouseHardware(houseID, request.Offset)
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, gin.H{
		"Hardware": hardware,
		"Count":    count,
	})
}

func handlerGetHardware(c *gin.Context) {
	request := struct {
		Offset int
	}{}

	if err := c.BindJSON(&request); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	hardware, count, err := database.GetHardware(request.Offset)
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, gin.H{
		"Hardware": hardware,
		"Count":    count,
	})
}

func handlerGetSwitches(c *gin.Context) {
	switches, err := database.GetSwitches()
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, switches)
}

func handlerEditSwitch(c *gin.Context) {
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

	if err := _switch.EditSwitch(); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, _switch)
}

func handlerCreateSwitch(c *gin.Context) {
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

	if err := _switch.CreateSwitch(); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, _switch)
}

func handlerGetHardwareTypes(c *gin.Context) {
	hardwareTypes, err := database.GetHardwareTypes()
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, hardwareTypes)
}

func handlerGetOperationModes(c *gin.Context) {
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

	operationModes, err := database.GetOperationModes()
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, operationModes)
}
