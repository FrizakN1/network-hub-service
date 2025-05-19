package models

import "database/sql"

type Hardware struct {
	ID          int
	Node        Node
	Type        Reference
	Switch      Switch
	IpAddress   sql.NullString
	MgmtVlan    sql.NullString
	Description sql.NullString
	CreatedAt   int64
	UpdatedAt   sql.NullInt64
	IsDelete    bool
}
