package handlers

import (
	"backend/database"
	"backend/errors"
	"backend/models"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type HardwareHandler interface {
	HandlerGetHardwareByID(c *gin.Context)
	HandlerEditHardware(c *gin.Context)
	HandlerCreateHardware(c *gin.Context)
	HandlerGetSearchHardware(c *gin.Context)
	HandlerGetNodeHardware(c *gin.Context)
	HandlerGetHouseHardware(c *gin.Context)
	HandlerGetHardware(c *gin.Context)
	HandlerDeleteHardware(c *gin.Context)
}

type DefaultHardwareHandler struct {
	Privilege    Privilege
	HardwareRepo database.HardwareRepository
	EventRepo    database.EventRepository
}

func NewHardwareHandler(db *database.Database) HardwareHandler {
	return &DefaultHardwareHandler{
		Privilege: &DefaultPrivilege{},
		HardwareRepo: &database.DefaultHardwareRepository{
			Database: *db,
			Counter: &database.DefaultCounter{
				Database: *db,
			},
		},
		EventRepo: &database.DefaultEventRepository{
			Database: *db,
			Counter: &database.DefaultCounter{
				Database: *db,
			},
		},
	}
}

func (h *DefaultHardwareHandler) HandlerDeleteHardware(c *gin.Context) {
	_, isAdmin, _ := h.Privilege.getPrivilege(c)

	if !isAdmin {
		c.Error(errors.NewHTTPError(nil, "forbidden", http.StatusForbidden))
		return
	}

	hardwareID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to parse param(id) to int", http.StatusBadRequest))
		return
	}

	if err = h.HardwareRepo.DeleteHardware(hardwareID); err != nil {
		c.Error(errors.NewHTTPError(err, "failed to delete hardware", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, true)
}

func (h *DefaultHardwareHandler) HandlerGetHardwareByID(c *gin.Context) {
	var (
		err      error
		hardware models.Hardware
	)

	hardware.ID, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to parse param(id)", http.StatusBadRequest))
		return
	}

	if err = h.HardwareRepo.GetHardwareByID(&hardware); err != nil {
		c.Error(errors.NewHTTPError(err, "failed to get hardware", http.StatusBadRequest))
		return
	}

	c.JSON(http.StatusOK, hardware)
}

func (h *DefaultHardwareHandler) HandlerEditHardware(c *gin.Context) {
	session, _, isOperatorOrHigher := h.Privilege.getPrivilege(c)

	if !isOperatorOrHigher {
		c.Error(errors.NewHTTPError(nil, "forbidden", http.StatusForbidden))
		return
	}

	var hardware models.Hardware

	if err := c.BindJSON(&hardware); err != nil {
		c.Error(errors.NewHTTPError(err, "invalid json", http.StatusBadRequest))
		return
	}

	if !h.HardwareRepo.ValidateHardware(hardware) {
		c.Error(errors.NewHTTPError(nil, "invalid hardware data", http.StatusBadRequest))
		return
	}

	hardware.UpdatedAt = sql.NullInt64{Int64: time.Now().Unix(), Valid: true}

	if err := h.HardwareRepo.EditHardware(&hardware); err != nil {
		c.Error(errors.NewHTTPError(err, "failed to edit hardware", http.StatusInternalServerError))
		return
	}

	event := models.Event{
		Address:     models.Address{House: models.AddressElement{ID: hardware.Node.Address.House.ID}},
		Node:        &models.Node{ID: hardware.Node.ID},
		Hardware:    &models.Hardware{ID: hardware.ID},
		UserId:      session.User.Id,
		Description: fmt.Sprintf("Изменение оборудования: %s", hardware.Type.TranslateValue),
		CreatedAt:   time.Now().Unix(),
	}

	if err := h.EventRepo.CreateEvent(event); err != nil {
		c.Error(errors.NewHTTPError(err, "failed to create event", http.StatusInternalServerError))
	}

	c.JSON(http.StatusOK, hardware)
}

func (h *DefaultHardwareHandler) HandlerCreateHardware(c *gin.Context) {
	session, _, isOperatorOrHigher := h.Privilege.getPrivilege(c)

	if !isOperatorOrHigher {
		c.Error(errors.NewHTTPError(nil, "forbidden", http.StatusForbidden))
		return
	}

	var hardware models.Hardware

	if err := c.BindJSON(&hardware); err != nil {
		c.Error(errors.NewHTTPError(err, "invalid json", http.StatusBadRequest))
		return
	}

	if !h.HardwareRepo.ValidateHardware(hardware) {
		c.Error(errors.NewHTTPError(nil, "invalid hardware data", http.StatusBadRequest))
		return
	}

	hardware.CreatedAt = time.Now().Unix()

	if err := h.HardwareRepo.CreateHardware(&hardware); err != nil {
		c.Error(errors.NewHTTPError(err, "failed to create hardware", http.StatusInternalServerError))
		return
	}

	event := models.Event{
		Address:     models.Address{House: models.AddressElement{ID: hardware.Node.Address.House.ID}},
		Node:        &models.Node{ID: hardware.Node.ID},
		Hardware:    nil,
		UserId:      session.User.Id,
		Description: fmt.Sprintf("Создание оборудования: %s", hardware.Type.TranslateValue),
		CreatedAt:   time.Now().Unix(),
	}

	if err := h.EventRepo.CreateEvent(event); err != nil {
		c.Error(errors.NewHTTPError(err, "failed to create event", http.StatusInternalServerError))
	}

	c.JSON(http.StatusOK, hardware)
}

func (h *DefaultHardwareHandler) HandlerGetSearchHardware(c *gin.Context) {
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to parse query(offset) to int", http.StatusBadRequest))
		return
	}
	search := c.Query("search")

	hardware, count, err := h.HardwareRepo.GetSearchHardware(search, offset)
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to get search hardware", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Hardware": hardware,
		"Count":    count,
	})
}

func (h *DefaultHardwareHandler) HandlerGetNodeHardware(c *gin.Context) {
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to parse query(offset) to int", http.StatusBadRequest))
		return
	}

	nodeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to parse param(id) to int", http.StatusBadRequest))
		return
	}

	hardware, count, err := h.HardwareRepo.GetNodeHardware(nodeID, offset)
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to get hardware", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Hardware": hardware,
		"Count":    count,
	})
}

func (h *DefaultHardwareHandler) HandlerGetHouseHardware(c *gin.Context) {
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to parse query(offset) to int", http.StatusBadRequest))
		return
	}

	houseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to parse param(id) to int", http.StatusBadRequest))
		return
	}

	hardware, count, err := h.HardwareRepo.GetHouseHardware(houseID, offset)
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to get hardware", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Hardware": hardware,
		"Count":    count,
	})
}

func (h *DefaultHardwareHandler) HandlerGetHardware(c *gin.Context) {
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to parse query(offset) to int", http.StatusBadRequest))
		return
	}

	hardware, count, err := h.HardwareRepo.GetHardware(offset)
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to get hardware", http.StatusBadRequest))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Hardware": hardware,
		"Count":    count,
	})
}
