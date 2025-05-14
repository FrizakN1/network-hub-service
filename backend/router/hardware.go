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

type HardwareHandler interface {
	handlerGetHardwareByID(c *gin.Context)
	handlerEditHardware(c *gin.Context)
	handlerCreateHardware(c *gin.Context)
	handlerGetSearchHardware(c *gin.Context)
	handlerGetNodeHardware(c *gin.Context)
	handlerGetHouseHardware(c *gin.Context)
	handlerGetHardware(c *gin.Context)
	handlerDeleteHardware(c *gin.Context)
}

type DefaultHardwareHandler struct {
	HardwareService database.HardwareService
}

func (hh *DefaultHardwareHandler) handlerDeleteHardware(c *gin.Context) {
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

	hardwareID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	if err = hh.HardwareService.DeleteHardware(hardwareID); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, true)
}

func (hh *DefaultHardwareHandler) handlerGetHardwareByID(c *gin.Context) {
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

	if err = hh.HardwareService.GetHardwareByID(&hardware); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, hardware)
}

func (hh *DefaultHardwareHandler) handlerEditHardware(c *gin.Context) {
	sessionHash, ok := c.Get("sessionHash")
	if !ok {
		err := errors.New("сессия не найдена")
		utils.Logger.Println(err)
		handlerError(c, err, 401)
		return
	}

	session := database.GetSession(sessionHash.(string))

	if session.User.Role.Value != "admin" && session.User.Role.Value != "operator" {
		c.JSON(403, nil)
		return
	}

	var hardware database.Hardware

	if err := c.BindJSON(&hardware); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	if !hh.HardwareService.ValidateHardware(hardware) {
		c.JSON(400, nil)
		return
	}

	hardware.UpdatedAt = sql.NullInt64{Int64: time.Now().Unix(), Valid: true}

	if err := hh.HardwareService.EditHardware(&hardware); err != nil {
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

func (hh *DefaultHardwareHandler) handlerCreateHardware(c *gin.Context) {
	sessionHash, ok := c.Get("sessionHash")
	if !ok {
		err := errors.New("сессия не найдена")
		utils.Logger.Println(err)
		handlerError(c, err, 401)
		return
	}

	session := database.GetSession(sessionHash.(string))

	if session.User.Role.Value != "admin" && session.User.Role.Value != "operator" {
		c.JSON(403, nil)
		return
	}

	var hardware database.Hardware

	if err := c.BindJSON(&hardware); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	if !hh.HardwareService.ValidateHardware(hardware) {
		c.JSON(400, nil)
		return
	}

	hardware.CreatedAt = time.Now().Unix()

	if err := hh.HardwareService.CreateHardware(&hardware); err != nil {
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

func (hh *DefaultHardwareHandler) handlerGetSearchHardware(c *gin.Context) {
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}
	search := c.Query("search")

	hardware, count, err := hh.HardwareService.GetSearchHardware(search, offset)
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

func (hh *DefaultHardwareHandler) handlerGetNodeHardware(c *gin.Context) {
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

	hardware, count, err := hh.HardwareService.GetNodeHardware(nodeID, offset)
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

func (hh *DefaultHardwareHandler) handlerGetHouseHardware(c *gin.Context) {
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

	hardware, count, err := hh.HardwareService.GetHouseHardware(houseID, offset)
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

func (hh *DefaultHardwareHandler) handlerGetHardware(c *gin.Context) {
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	hardware, count, err := hh.HardwareService.GetHardware(offset)
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

func NewHardwareHandler() HardwareHandler {
	return &DefaultHardwareHandler{
		HardwareService: &database.DefaultHardwareService{},
	}
}
