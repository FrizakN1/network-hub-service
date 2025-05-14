package router

import (
	"backend/database"
	"backend/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

type ReferenceHandler interface {
	handleReferenceRecord(c *gin.Context, isEdit bool)
	handlerGetReference(c *gin.Context, onlyAdmin bool)
}

type DefaultReferenceHandler struct {
	ReferenceService database.ReferenceService
}

func (rh *DefaultReferenceHandler) handleReferenceRecord(c *gin.Context, isEdit bool) {
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

	var record database.Reference
	reference := c.Param("reference")

	if err := c.BindJSON(&record); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	if ((reference == "node_types" || reference == "owners") && record.Name == "") ||
		((reference == "hardware_types" || reference == "operation_modes") && record.Value == "" && record.TranslateValue == "") {
		c.JSON(400, nil)
		return
	}

	if !isEdit {
		record.CreatedAt = time.Now().Unix()
		err := rh.ReferenceService.CreateReferenceRecord(&record, strings.ToUpper(reference))
		if err != nil {
			utils.Logger.Println(err)
			handlerError(c, err, 400)
			return
		}
	} else {
		err := rh.ReferenceService.EditReferenceRecord(&record, strings.ToUpper(reference))
		if err != nil {
			utils.Logger.Println(err)
			handlerError(c, err, 400)
			return
		}
	}

	c.JSON(200, record)
}

func (rh *DefaultReferenceHandler) handlerGetReference(c *gin.Context, onlyAdmin bool) {
	if onlyAdmin {
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
	}

	reference := c.Param("reference")

	records, err := rh.ReferenceService.GetReferenceRecords(strings.ToUpper(reference))
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, records)
}

func NewReferenceHandler() ReferenceHandler {
	return &DefaultReferenceHandler{
		ReferenceService: &database.DefaultReferenceService{},
	}
}
