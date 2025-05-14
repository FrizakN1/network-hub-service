package router

import (
	"backend/database"
	"backend/proto/userpb"
	"backend/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

func (h *DefaultHandler) handleReferenceRecord(c *gin.Context, isEdit bool) {
	session, ok := c.Get("session")
	if !ok {
		err := errors.New("сессия не найдена")
		utils.Logger.Println(err)
		handlerError(c, err, 401)
		return
	}

	if session.(userpb.Session).User.Role.Value != "admin" && session.(userpb.Session).User.Role.Value != "operator" {
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
		err := h.ReferenceService.CreateReferenceRecord(&record, strings.ToUpper(reference))
		if err != nil {
			utils.Logger.Println(err)
			handlerError(c, err, 400)
			return
		}
	} else {
		err := h.ReferenceService.EditReferenceRecord(&record, strings.ToUpper(reference))
		if err != nil {
			utils.Logger.Println(err)
			handlerError(c, err, 400)
			return
		}
	}

	c.JSON(200, record)
}

func (h *DefaultHandler) handlerGetReference(c *gin.Context) {
	reference := c.Param("reference")

	records, err := h.ReferenceService.GetReferenceRecords(strings.ToUpper(reference))
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, records)
}
