package models

import "database/sql"

type Report struct {
	ID          int
	Key         string
	Value       string
	Description sql.NullString
}
