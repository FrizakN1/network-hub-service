package database

import (
	"backend/utils"
	"database/sql"
	"errors"
)

type Reference struct {
	ID             int
	Name           string
	Value          string
	TranslateValue string
	CreatedAt      int64
}

var query map[string]*sql.Stmt

func prepareReferences() []string {
	var e error
	errorsList := make([]string, 0)

	if query == nil {
		query = make(map[string]*sql.Stmt)
	}

	query["GET_ROLES"], e = Link.Prepare(`
		SELECT * FROM "Role"
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["GET_OWNERS"], e = Link.Prepare(`
		SELECT * FROM "Node_owner"
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["GET_NODE_TYPES"], e = Link.Prepare(`
		SELECT * FROM "Node_type"
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["GET_HARDWARE_TYPES"], e = Link.Prepare(`
		SELECT * FROM "Hardware_type" ORDER BY id
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["GET_OPERATION_MODES"], e = Link.Prepare(`
		SELECT * FROM "Operation_mode" ORDER BY id
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["CREATE_OWNERS"], e = Link.Prepare(`
		INSERT INTO "Node_owner"(name, created_at) VALUES ($1, $2)
		RETURNING id
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["EDIT_OWNERS"], e = Link.Prepare(`
		UPDATE "Node_owner" SET name = $2 WHERE id = $1
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["CREATE_NODE_TYPES"], e = Link.Prepare(`
		INSERT INTO "Node_type"(name, created_at) VALUES ($1, $2)
		RETURNING id
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["EDIT_NODE_TYPES"], e = Link.Prepare(`
		UPDATE "Node_type" SET name = $2 WHERE id = $1
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["CREATE_HARDWARE_TYPES"], e = Link.Prepare(`
		INSERT INTO "Hardware_type"(value, translate_value, created_at) VALUES ($1, $2, $3)
		RETURNING id
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["EDIT_HARDWARE_TYPES"], e = Link.Prepare(`
		UPDATE "Hardware_type" SET value = $2, translate_value = $3 WHERE id = $1
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["CREATE_OPERATION_MODES"], e = Link.Prepare(`
		INSERT INTO "Operation_mode"(value, translate_value, created_at) VALUES ($1, $2, $3)
		RETURNING id
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["EDIT_OPERATION_MODES"], e = Link.Prepare(`
		UPDATE "Operation_mode" SET value = $2, translate_value = $3 WHERE id = $1
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	return errorsList
}

func (referenceRecord *Reference) EditReferenceRecord(reference string) error {
	stmt, ok := query["EDIT_"+reference]
	if !ok {
		err := errors.New("запрос EDIT_" + reference + " не подготовлен")
		utils.Logger.Println(err)
		return err
	}

	var err error

	switch reference {
	case "NODE_TYPES":
	case "OWNERS":
		_, err = stmt.Exec(referenceRecord.ID, referenceRecord.Name)
		break
	case "HARDWARE_TYPES":
	case "OPERATION_MODES":
		_, err = stmt.Exec(referenceRecord.ID, referenceRecord.Value, referenceRecord.TranslateValue)
		break
	}

	if err != nil {
		utils.Logger.Println(err)
		return err
	}

	return nil
}

func (referenceRecord *Reference) CreateReferenceRecord(reference string) error {
	stmt, ok := query["CREATE_"+reference]
	if !ok {
		err := errors.New("запрос CREATE_" + reference + " не подготовлен")
		utils.Logger.Println(err)
		return err
	}

	if reference == "NODE_TYPES" || reference == "OWNERS" {
		if err := stmt.QueryRow(
			referenceRecord.Name,
			referenceRecord.CreatedAt,
		).Scan(&referenceRecord.ID); err != nil {
			utils.Logger.Println(err)
			return err
		}
	} else {
		if err := stmt.QueryRow(
			referenceRecord.Value,
			referenceRecord.TranslateValue,
			referenceRecord.CreatedAt,
		).Scan(&referenceRecord.ID); err != nil {
			utils.Logger.Println(err)
			return err
		}
	}

	return nil
}

func GetReferenceRecords(reference string) ([]Reference, error) {
	stmt, ok := query["GET_"+reference]
	if !ok {
		err := errors.New("запрос GET_" + reference + " не подготовлен")
		utils.Logger.Println(err)
		return nil, err
	}

	rows, err := stmt.Query()
	if err != nil {
		utils.Logger.Println(err)
		return nil, err
	}
	defer rows.Close()

	var references []Reference
	for rows.Next() {
		var _reference Reference

		if reference == "NODE_TYPES" || reference == "OWNERS" {
			err = rows.Scan(&_reference.ID, &_reference.Name, &_reference.CreatedAt)
		} else if reference == "HARDWARE_TYPES" || reference == "OPERATION_MODES" {
			err = rows.Scan(&_reference.ID, &_reference.Value, &_reference.TranslateValue, &_reference.CreatedAt)
		} else if reference == "ROLES" {
			err = rows.Scan(&_reference.ID, &_reference.Value, &_reference.TranslateValue)
		}

		if err != nil {
			utils.Logger.Println(err)
			return nil, err
		}

		references = append(references, _reference)
	}

	return references, nil
}
