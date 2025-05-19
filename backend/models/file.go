package models

type File struct {
	ID             int
	House          AddressElement
	Node           Node
	Hardware       Hardware
	Path           string
	Name           string
	UploadAt       int64
	Data           string
	InArchive      bool
	IsPreviewImage bool
}
