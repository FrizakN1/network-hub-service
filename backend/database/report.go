package database

import (
	"backend/models"
	"errors"
)

type ReportRepository interface {
	GetReportData() ([]models.Report, error)
	EditReportData(reportData *models.Report) error
}

type DefaultReportRepository struct {
	Database Database
}

func (r *DefaultReportRepository) EditReportData(reportData *models.Report) error {
	stmt, ok := r.Database.GetQuery("EDIT_REPORT_DATA")
	if !ok {
		return errors.New("query EDIT_REPORT_DATA is not prepare")
	}

	_, err := stmt.Exec(reportData.Key, reportData.Value, reportData.Description)
	if err != nil {
		return err
	}

	return nil
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
