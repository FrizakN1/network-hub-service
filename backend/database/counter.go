package database

import (
	"database/sql"
	"errors"
)

type Counter interface {
	countRecords(key string, param []interface{}) (int, error)
}

type DefaultCounter struct {
	Database Database
}

func (c *DefaultCounter) countRecords(key string, param []interface{}) (int, error) {
	stmt, ok := c.Database.GetQuery(key)
	if !ok {
		return 0, errors.New("запрос " + key + " не подготовлен")
	}

	var row *sql.Row

	if param != nil {
		row = stmt.QueryRow(param...)
	} else {
		row = stmt.QueryRow()
	}

	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

//func (r *DefaultAddressRepository) countSuggestions(streetPart, housePart string) (int, error) {
//	stmt, ok := r.Database.GetQuery("GET_SUGGESTIONS_COUNT")
//	if !ok {
//		return 0, errors.New("запрос GET_SUGGESTIONS_COUNT не подготовлен")
//	}
//
//	row := stmt.QueryRow(streetPart, housePart)
//
//	var count int
//	err := row.Scan(&count)
//	if err != nil {
//		return 0, err
//	}
//
//	return count, nil
//}
