package models

type AddressParams struct {
	HouseID        int
	FileAmount     int
	NodeAmount     int
	HardwareAmount int
	RoofType       Reference
	WiringType     Reference
}
