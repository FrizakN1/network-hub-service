package handlers

import (
	"backend/database"
	"backend/errors"
	"backend/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type SwitchHandler interface {
	HandlerGetSwitches(c *gin.Context)
	HandlerEditSwitch(c *gin.Context)
	HandlerCreateSwitch(c *gin.Context)
}

type DefaultSwitchHandler struct {
	Privilege  Privilege
	SwitchRepo database.SwitchRepository
}

func NewSwitchHandler(db *database.Database) SwitchHandler {
	return &DefaultSwitchHandler{
		Privilege: &DefaultPrivilege{},
		SwitchRepo: &database.DefaultSwitchRepository{
			Database: *db,
		},
	}
}

func (h *DefaultSwitchHandler) HandlerGetSwitches(c *gin.Context) {
	switches, err := h.SwitchRepo.GetSwitches()
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to get switches", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, switches)
}

func (h *DefaultSwitchHandler) HandlerEditSwitch(c *gin.Context) {
	_, _, isOperatorOrHigher := h.Privilege.getPrivilege(c)

	if !isOperatorOrHigher {
		c.Error(errors.NewHTTPError(nil, "forbidden", http.StatusForbidden))
		return
	}

	var _switch models.Switch

	if err := c.BindJSON(&_switch); err != nil {
		c.Error(errors.NewHTTPError(err, "invalid json", http.StatusBadRequest))
		return
	}

	if len(_switch.Name) == 0 || _switch.PortAmount == 0 {
		c.Error(errors.NewHTTPError(nil, "invalid switch data", http.StatusBadRequest))
		return
	}

	if err := h.SwitchRepo.EditSwitch(&_switch); err != nil {
		c.Error(errors.NewHTTPError(err, "failed to edit switch", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, _switch)
}

func (h *DefaultSwitchHandler) HandlerCreateSwitch(c *gin.Context) {
	_, _, isOperatorOrHigher := h.Privilege.getPrivilege(c)

	if !isOperatorOrHigher {
		c.Error(errors.NewHTTPError(nil, "forbidden", http.StatusForbidden))
		return
	}

	var _switch models.Switch

	if err := c.BindJSON(&_switch); err != nil {
		c.Error(errors.NewHTTPError(err, "invalid json", http.StatusBadRequest))
		return
	}

	if len(_switch.Name) == 0 || _switch.PortAmount == 0 {
		c.Error(errors.NewHTTPError(nil, "invalid switch data", http.StatusBadRequest))
		return
	}

	_switch.CreatedAt = time.Now().Unix()

	if err := h.SwitchRepo.CreateSwitch(&_switch); err != nil {
		c.Error(errors.NewHTTPError(err, "failed to create switch", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, _switch)
}
