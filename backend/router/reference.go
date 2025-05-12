package router

import (
	"backend/database"
	"backend/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

func handleReferenceRecord(c *gin.Context, isEdit bool) {
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
		err := record.CreateReferenceRecord(strings.ToUpper(reference))
		if err != nil {
			utils.Logger.Println(err)
			handlerError(c, err, 400)
			return
		}
	} else {
		err := record.EditReferenceRecord(strings.ToUpper(reference))
		if err != nil {
			utils.Logger.Println(err)
			handlerError(c, err, 400)
			return
		}
	}

	c.JSON(200, record)
}

func handlerGetReference(c *gin.Context, onlyAdmin bool) {
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

	records, err := database.GetReferenceRecords(strings.ToUpper(reference))
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, records)
}
