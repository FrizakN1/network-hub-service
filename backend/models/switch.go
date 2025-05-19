package models

import "database/sql"

type Switch struct {
	ID               int
	Name             string
	OperationMode    Reference
	PortAmount       int
	CommunityRead    sql.NullString
	CommunityWrite   sql.NullString
	FirmwareOID      sql.NullString
	SystemNameOID    sql.NullString
	SerialNumberOID  sql.NullString
	SaveConfigOID    sql.NullString
	PortDescOID      sql.NullString
	VlanOID          sql.NullString
	PortUntaggedOID  sql.NullString
	SpeedOID         sql.NullString
	BatteryStatusOID sql.NullString
	BatteryChargeOID sql.NullString
	PortModeOID      sql.NullString
	UptimeOID        sql.NullString
	CreatedAt        int64
	MacOID           sql.NullString
}
