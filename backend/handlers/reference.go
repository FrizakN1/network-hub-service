package handlers

import (
	"backend/database"
	"backend/errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

type ReferenceHandler interface {
	HandlerReferenceRecord(c *gin.Context, isEdit bool)
	HandlerGetReference(c *gin.Context)
}

type DefaultReferenceHandler struct {
	Privilege        Privilege
	ReferenceService database.ReferenceService
}

func NewReferenceHandler() ReferenceHandler {
	return &DefaultReferenceHandler{
		Privilege:        &DefaultPrivilege{},
		ReferenceService: &database.DefaultReferenceService{},
	}
}

func (h *DefaultReferenceHandler) HandlerReferenceRecord(c *gin.Context, isEdit bool) {
	_, _, isOperatorOrHigher := h.Privilege.getPrivilege(c)

	if !isOperatorOrHigher {
		c.Error(errors.NewHTTPError(nil, "forbidden", http.StatusForbidden))
		return
	}

	var record database.Reference
	reference := c.Param("reference")

	if err := c.BindJSON(&record); err != nil {
		c.Error(errors.NewHTTPError(nil, "invalid json data", http.StatusBadRequest))
		return
	}

	if ((reference == "node_types" || reference == "owners") && record.Name == "") ||
		((reference == "hardware_types" || reference == "operation_modes") && record.Value == "" && record.TranslateValue == "") {
		c.Error(errors.NewHTTPError(nil, fmt.Sprintf("invalid %s data", reference), http.StatusBadRequest))
		return
	}

	if !isEdit {
		record.CreatedAt = time.Now().Unix()
		err := h.ReferenceService.CreateReferenceRecord(&record, strings.ToUpper(reference))
		if err != nil {
			c.Error(errors.NewHTTPError(err, fmt.Sprintf("failed to create %s", reference), http.StatusInternalServerError))
			return
		}
	} else {
		err := h.ReferenceService.EditReferenceRecord(&record, strings.ToUpper(reference))
		if err != nil {
			c.Error(errors.NewHTTPError(err, "failed to edit reference", http.StatusInternalServerError))
			return
		}
	}

	c.JSON(http.StatusOK, record)
}

func (h *DefaultReferenceHandler) HandlerGetReference(c *gin.Context) {
	reference := c.Param("reference")

	records, err := h.ReferenceService.GetReferenceRecords(strings.ToUpper(reference))
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to get references", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, records)
}
