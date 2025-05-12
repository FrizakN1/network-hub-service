package router

import (
	"backend/database"
	"backend/utils"
	"database/sql"
	"errors"
	"fmt"
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
	sessionHash, ok := c.Get("sessionHash")
	if !ok {
		err := errors.New("сессия не найдена")
		utils.Logger.Println(err)
		handlerError(c, err, 401)
		return
	}

	session := database.GetSession(sessionHash.(string))

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

	event := database.Event{
		Address:     database.Address{House: database.AddressElement{ID: hardware.Node.Address.House.ID}},
		Node:        &database.Node{ID: hardware.Node.ID},
		Hardware:    &database.Hardware{ID: hardware.ID},
		User:        database.User{ID: session.User.ID},
		Description: fmt.Sprintf("Изменение оборудования: %s", hardware.Type.TranslateValue),
		CreatedAt:   time.Now().Unix(),
	}

	if err := event.CreateEvent(); err != nil {
		utils.Logger.Println(err)
	}

	c.JSON(200, hardware)
}

func handlerCreateHardware(c *gin.Context) {
	sessionHash, ok := c.Get("sessionHash")
	if !ok {
		err := errors.New("сессия не найдена")
		utils.Logger.Println(err)
		handlerError(c, err, 401)
		return
	}

	session := database.GetSession(sessionHash.(string))

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

	event := database.Event{
		Address:     database.Address{House: database.AddressElement{ID: hardware.Node.Address.House.ID}},
		Node:        &database.Node{ID: hardware.Node.ID},
		Hardware:    nil,
		User:        database.User{ID: session.User.ID},
		Description: fmt.Sprintf("Создание оборудования: %s", hardware.Type.TranslateValue),
		CreatedAt:   time.Now().Unix(),
	}

	if err := event.CreateEvent(); err != nil {
		utils.Logger.Println(err)
	}

	c.JSON(200, hardware)
}

func handlerGetSearchHardware(c *gin.Context) {
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}
	search := c.Query("search")

	hardware, count, err := database.GetSearchHardware(search, offset)
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
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
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

	hardware, count, err := database.GetNodeHardware(nodeID, offset)
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
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
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

	hardware, count, err := database.GetHouseHardware(houseID, offset)
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
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	hardware, count, err := database.GetHardware(offset)
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
