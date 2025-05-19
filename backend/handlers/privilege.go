package handlers

import (
	"backend/proto/userpb"
	"github.com/gin-gonic/gin"
)

type Privilege interface {
	getPrivilege(c *gin.Context) (*userpb.Session, bool, bool)
}

type DefaultPrivilege struct{}

func (p *DefaultPrivilege) getPrivilege(c *gin.Context) (*userpb.Session, bool, bool) {
	session, ok := c.Get("session")
	if !ok {
		return nil, false, false
	}

	isAdmin := session.(*userpb.Session).User.Role.Key == "admin"
	isOperatorOrHigher := session.(*userpb.Session).User.Role.Key == "operator" || isAdmin

	return session.(*userpb.Session), isAdmin, isOperatorOrHigher
}
