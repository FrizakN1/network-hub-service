package models

import (
	"backend/proto/addresspb"
	"backend/proto/userpb"
)

type Event struct {
	ID          int64
	HouseId     int32
	Address     *addresspb.Address
	Node        *Node
	Hardware    *Hardware
	UserId      int32
	User        *userpb.User
	Description string
	CreatedAt   int64
}
