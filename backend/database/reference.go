package database

import (
	"backend/models"
	"errors"
)

type ReferenceRepository interface {
	EditReferenceRecord(referenceRecord *models.Reference, reference string) error
	CreateReferenceRecord(referenceRecord *models.Reference, reference string) error
	GetReferenceRecords(reference string) ([]models.Reference, error)
}

type DefaultReferenceRepository struct {
	Database Database
}

func (r *DefaultReferenceRepository) EditReferenceRecord(referenceRecord *models.Reference, reference string) error {
	stmt, ok := r.Database.GetQuery("EDIT_" + reference)
	if !ok {
		return errors.New("запрос EDIT_" + reference + " не подготовлен")
	}

	var err error

	if reference == "NODE_TYPES" || reference == "OWNERS" {
		_, err = stmt.Exec(referenceRecord.ID, referenceRecord.Name)
	} else if reference == "HARDWARE_TYPES" || reference == "OPERATION_MODES" {
		_, err = stmt.Exec(referenceRecord.ID, referenceRecord.Value, referenceRecord.TranslateValue)
	}

	if err != nil {
		return err
	}

	return nil
}

func (r *DefaultReferenceRepository) CreateReferenceRecord(referenceRecord *models.Reference, reference string) error {
	stmt, ok := r.Database.GetQuery("CREATE_" + reference)
	if !ok {
		return errors.New("запрос CREATE_" + reference + " не подготовлен")
	}

	if reference == "NODE_TYPES" || reference == "OWNERS" {
		if err := stmt.QueryRow(
			referenceRecord.Name,
			referenceRecord.CreatedAt,
		).Scan(&referenceRecord.ID); err != nil {
			return err
		}
	} else {
		if err := stmt.QueryRow(
			referenceRecord.Value,
			referenceRecord.TranslateValue,
			referenceRecord.CreatedAt,
		).Scan(&referenceRecord.ID); err != nil {
			return err
		}
	}

	return nil
}

func (r *DefaultReferenceRepository) GetReferenceRecords(reference string) ([]models.Reference, error) {
	stmt, ok := r.Database.GetQuery("GET_" + reference)
	if !ok {
		return nil, errors.New("запрос GET_" + reference + " не подготовлен")
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var references []models.Reference
	for rows.Next() {
		var _reference models.Reference

		if reference == "NODE_TYPES" || reference == "OWNERS" {
			err = rows.Scan(&_reference.ID, &_reference.Name, &_reference.CreatedAt)
		} else if reference == "HARDWARE_TYPES" || reference == "OPERATION_MODES" {
			err = rows.Scan(&_reference.ID, &_reference.Value, &_reference.TranslateValue, &_reference.CreatedAt)
		} else if reference == "ROLES" {
			err = rows.Scan(&_reference.ID, &_reference.Value, &_reference.TranslateValue)
		}

		if err != nil {
			return nil, err
		}

		references = append(references, _reference)
	}

	return references, nil
}
