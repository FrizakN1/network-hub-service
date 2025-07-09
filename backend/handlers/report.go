package handlers

import (
	"backend/database"
	"backend/errors"
	"backend/models"
	"backend/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ReportHandler interface {
	HandlerGetReportData(c *gin.Context)
	HandlerEditReportData(c *gin.Context)
}

type DefaultReportHandler struct {
	Privilege  Privilege
	ReportRepo database.ReportRepository
	utils.Logger
}

func NewReportHandler(db *database.Database, logger *utils.Logger) ReportHandler {
	return &DefaultReportHandler{
		Privilege: &DefaultPrivilege{},

		ReportRepo: &database.DefaultReportRepository{
			Database: *db,
		},
		Logger: *logger,
	}
}

func (h *DefaultReportHandler) HandlerEditReportData(c *gin.Context) {
	_, _, isOperatorOrHigher := h.Privilege.getPrivilege(c)

	if !isOperatorOrHigher {
		c.Error(errors.NewHTTPError(nil, "forbidden", http.StatusForbidden))
		return
	}

	reportData := &models.Report{}

	if err := c.BindJSON(reportData); err != nil {
		c.Error(errors.NewHTTPError(err, "invalid json", http.StatusBadRequest))
		return
	}

	if err := h.ReportRepo.EditReportData(reportData); err != nil {
		c.Error(errors.NewHTTPError(err, "failed to edit report data", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, reportData)
}

func (h *DefaultReportHandler) HandlerGetReportData(c *gin.Context) {
	reportData, err := h.ReportRepo.GetReportData()
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to get report data", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, reportData)
}
