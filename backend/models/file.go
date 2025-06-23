package models

import "backend/proto/addresspb"

type File struct {
	ID             int
	HouseId        int32
	Address        *addresspb.Address
	Node           Node
	Hardware       Hardware
	Path           string
	Name           string
	UploadAt       int64
	Data           string
	InArchive      bool
	IsPreviewImage bool
}
