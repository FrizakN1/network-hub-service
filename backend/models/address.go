package models

import "database/sql"

type AddressElementType struct {
	ID        int
	Name      string
	ShortName string
}

type AddressElement struct {
	ID   int
	Name string
	Type AddressElementType
	FIAS sql.NullString
}

type Address struct {
	Street         AddressElement
	House          AddressElement
	FileAmount     int
	NodeAmount     int
	HardwareAmount int
	RoofType       Reference
	WiringType     Reference
}
