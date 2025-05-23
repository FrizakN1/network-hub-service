package database

import (
	"backend/models"
	"errors"
	"fmt"
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
		return errors.New("query EDIT_" + reference + " is not prepare")
	}

	var params []interface{}

	switch reference {
	case "NODE_TYPES", "OWNERS", "ROOF_TYPE", "WIRING_TYPE":
		params = []interface{}{referenceRecord.ID, referenceRecord.Value}
	case "HARDWARE_TYPES", "OPERATION_MODES":
		params = []interface{}{referenceRecord.ID, referenceRecord.Key, referenceRecord.Value}
	default:
		return fmt.Errorf("reference is unsupported (%s)", reference)
	}

	_, err := stmt.Exec(params...)
	if err != nil {
		return err
	}

	return nil
}

func (r *DefaultReferenceRepository) CreateReferenceRecord(referenceRecord *models.Reference, reference string) error {
	stmt, ok := r.Database.GetQuery("CREATE_" + reference)
	if !ok {
		return errors.New("query CREATE_" + reference + " is not prepare")
	}

	var params []interface{}

	switch reference {
	case "NODE_TYPES", "OWNERS", "ROOF_TYPES", "WIRING_TYPES":
		params = []interface{}{referenceRecord.Value, referenceRecord.CreatedAt}
	case "HARDWARE_TYPES", "OPERATION_MODES":
		params = []interface{}{referenceRecord.Key, referenceRecord.Value, referenceRecord.CreatedAt}
	default:
		return fmt.Errorf("reference is unsupported (%s)", reference)
	}

	if err := stmt.QueryRow(params...).Scan(&referenceRecord.ID); err != nil {
		return err
	}

	return nil
}

func (r *DefaultReferenceRepository) GetReferenceRecords(reference string) ([]models.Reference, error) {
	stmt, ok := r.Database.GetQuery("GET_" + reference)
	if !ok {
		return nil, errors.New("query GET_" + reference + " is not prepare")
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var references []models.Reference

	for rows.Next() {
		var ref models.Reference

		switch reference {
		case "NODE_TYPES", "OWNERS", "ROOF_TYPES", "WIRING_TYPES":
			err = rows.Scan(&ref.ID, &ref.Value, &ref.CreatedAt)
		case "HARDWARE_TYPES", "OPERATION_MODES":
			err = rows.Scan(&ref.ID, &ref.Key, &ref.Value, &ref.CreatedAt)
		default:
			return nil, fmt.Errorf("reference is unsupported (%s)", reference)
		}

		if err != nil {
			return nil, err
		}

		references = append(references, ref)
	}

	return references, nil
}
