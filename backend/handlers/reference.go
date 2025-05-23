package handlers

import (
	"backend/database"
	"backend/errors"
	"backend/models"
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
	Privilege     Privilege
	ReferenceRepo database.ReferenceRepository
}

func NewReferenceHandler(db *database.Database) ReferenceHandler {
	return &DefaultReferenceHandler{
		Privilege: &DefaultPrivilege{},
		ReferenceRepo: &database.DefaultReferenceRepository{
			Database: *db,
		},
	}
}

func (h *DefaultReferenceHandler) HandlerReferenceRecord(c *gin.Context, isEdit bool) {
	_, _, isOperatorOrHigher := h.Privilege.getPrivilege(c)

	if !isOperatorOrHigher {
		c.Error(errors.NewHTTPError(nil, "forbidden", http.StatusForbidden))
		return
	}

	var record models.Reference
	reference := c.Param("reference")

	if err := c.BindJSON(&record); err != nil {
		c.Error(errors.NewHTTPError(nil, "invalid json data", http.StatusBadRequest))
		return
	}

	var isValid bool

	switch reference {
	case "node_types", "owners", "roof_types", "wiring_types":
		isValid = record.Value != ""
	case "hardware_types", "operation_modes":
		isValid = record.Key != "" && record.Value != ""
	default:
		c.Error(errors.NewHTTPError(nil, fmt.Sprintf("reference is unsupported (%s)", reference), http.StatusBadRequest))
		return
	}

	if !isValid {
		c.Error(errors.NewHTTPError(nil, fmt.Sprintf("invalid %s data", reference), http.StatusBadRequest))
		return
	}

	if !isEdit {
		record.CreatedAt = time.Now().Unix()
		err := h.ReferenceRepo.CreateReferenceRecord(&record, strings.ToUpper(reference))
		if err != nil {
			c.Error(errors.NewHTTPError(err, fmt.Sprintf("failed to create %s", reference), http.StatusInternalServerError))
			return
		}
	} else {
		err := h.ReferenceRepo.EditReferenceRecord(&record, strings.ToUpper(reference))
		if err != nil {
			c.Error(errors.NewHTTPError(err, "failed to edit reference", http.StatusInternalServerError))
			return
		}
	}

	c.JSON(http.StatusOK, record)
}

func (h *DefaultReferenceHandler) HandlerGetReference(c *gin.Context) {
	reference := c.Param("reference")

	records, err := h.ReferenceRepo.GetReferenceRecords(strings.ToUpper(reference))
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to get references", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, records)
}
