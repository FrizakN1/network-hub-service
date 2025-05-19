package models

import "backend/proto/userpb"

type Event struct {
	ID          int64
	Address     Address
	Node        *Node
	Hardware    *Hardware
	UserId      int32
	User        *userpb.User
	Description string
	CreatedAt   int64
}
