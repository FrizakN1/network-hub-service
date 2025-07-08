package database

import (
	"backend/models"
	"errors"
)

type ReportRepository interface {
	GetReportData() ([]models.Report, error)
}

type DefaultReportRepository struct {
	Database Database
}

func (r *DefaultReportRepository) GetReportData() ([]models.Report, error) {
	stmt, ok := r.Database.GetQuery("GET_REPORT_DATA")
	if !ok {
		return nil, errors.New("query GET_REPORT_DATA is not prepare")
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reportData []models.Report

	for rows.Next() {
		var report models.Report

		if err = rows.Scan(
			&report.ID,
			&report.Key,
			&report.Value,
			&report.Description,
		); err != nil {
			return nil, err
		}

		reportData = append(reportData, report)
	}

	return reportData, nil
}
