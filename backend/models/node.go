package models

import "database/sql"

type Node struct {
	ID          int
	Parent      *Node
	Address     Address
	Type        Reference
	Owner       Reference
	Name        string
	Zone        sql.NullString
	Placement   sql.NullString
	Supply      sql.NullString
	Access      sql.NullString
	Description sql.NullString
	CreatedAt   int64
	UpdatedAt   sql.NullInt64
	IsDelete    bool
	IsPassive   bool
}
