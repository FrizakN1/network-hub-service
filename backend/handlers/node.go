package handlers

import (
	"backend/database"
	"backend/errors"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type NodeHandler interface {
	HandlerGetSearchNodes(c *gin.Context)
	HandlerEditNode(c *gin.Context)
	HandlerCreateNode(c *gin.Context)
	HandlerGetNode(c *gin.Context)
	HandlerGetHouseNodes(c *gin.Context)
	HandlerGetNodes(c *gin.Context)
	HandlerDeleteNode(c *gin.Context)
}

type DefaultNodeHandler struct {
	Privilege    Privilege
	NodeService  database.NodeService
	EventService database.EventService
}

func NewNodeHandler() NodeHandler {
	return &DefaultNodeHandler{
		Privilege:    &DefaultPrivilege{},
		NodeService:  &database.DefaultNodeService{},
		EventService: &database.DefaultEventService{},
	}
}

func (h *DefaultNodeHandler) HandlerDeleteNode(c *gin.Context) {
	_, isAdmin, _ := h.Privilege.getPrivilege(c)

	if !isAdmin {
		c.Error(errors.NewHTTPError(nil, "forbidden", http.StatusForbidden))
		return
	}

	nodeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to parse param(id) to int", http.StatusBadRequest))
		return
	}

	if err = h.NodeService.DeleteNode(nodeID); err != nil {
		c.Error(errors.NewHTTPError(err, "failed to delete node", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, true)
}

func (h *DefaultNodeHandler) HandlerGetSearchNodes(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	search := c.Query("search")

	nodes, count, err := h.NodeService.GetSearchNodes(search, offset)
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to get search nodes", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Nodes": nodes,
		"Count": count,
	})
}

func (h *DefaultNodeHandler) HandlerEditNode(c *gin.Context) {
	session, _, isOperatorOrHigher := h.Privilege.getPrivilege(c)

	if !isOperatorOrHigher {
		c.Error(errors.NewHTTPError(nil, "forbidden", http.StatusForbidden))
		return
	}

	var node database.Node

	if err := c.BindJSON(&node); err != nil {
		c.Error(errors.NewHTTPError(err, "invalid json", http.StatusBadRequest))
		return
	}

	if !h.NodeService.ValidateNode(node) {
		c.Error(errors.NewHTTPError(nil, "invalid node data", http.StatusBadRequest))
		return
	}

	node.UpdatedAt = sql.NullInt64{
		Int64: time.Now().Unix(),
		Valid: true,
	}

	if err := h.NodeService.EditNode(&node); err != nil {
		c.Error(errors.NewHTTPError(err, "failed to edit node", http.StatusInternalServerError))
		return
	}

	event := database.Event{
		Address:     database.Address{House: database.AddressElement{ID: node.Address.House.ID}},
		Node:        &database.Node{ID: node.ID},
		Hardware:    nil,
		UserId:      session.User.Id,
		Description: fmt.Sprintf("Изменение узла: %s", node.Name),
		CreatedAt:   time.Now().Unix(),
	}

	if err := h.EventService.CreateEvent(event); err != nil {
		c.Error(errors.NewHTTPError(err, "failed to delete node", http.StatusInternalServerError))
	}

	c.JSON(http.StatusOK, node)
}

func (h *DefaultNodeHandler) HandlerCreateNode(c *gin.Context) {
	session, _, isOperatorOrHigher := h.Privilege.getPrivilege(c)

	if !isOperatorOrHigher {
		c.Error(errors.NewHTTPError(nil, "forbidden", http.StatusForbidden))
		return
	}

	var node database.Node

	if err := c.BindJSON(&node); err != nil {
		c.Error(errors.NewHTTPError(err, "invalid json", http.StatusBadRequest))
		return
	}

	if !h.NodeService.ValidateNode(node) {
		c.Error(errors.NewHTTPError(nil, "invalid node data", http.StatusBadRequest))
		return
	}

	node.CreatedAt = time.Now().Unix()

	if err := h.NodeService.CreateNode(&node); err != nil {
		c.Error(errors.NewHTTPError(err, "failed to create node", http.StatusInternalServerError))
		return
	}

	event := database.Event{
		Address:     database.Address{House: database.AddressElement{ID: node.Address.House.ID}},
		Node:        nil,
		Hardware:    nil,
		UserId:      session.User.Id,
		Description: fmt.Sprintf("Создание нового узла: %s", node.Name),
		CreatedAt:   time.Now().Unix(),
	}

	if err := h.EventService.CreateEvent(event); err != nil {
		c.Error(errors.NewHTTPError(err, "failed to create event", http.StatusInternalServerError))
	}

	c.JSON(http.StatusOK, node)
}

func (h *DefaultNodeHandler) HandlerGetNode(c *gin.Context) {
	var (
		err  error
		node database.Node
	)

	node.ID, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to parse param(id) to int", http.StatusBadRequest))
		return
	}

	if err = h.NodeService.GetNode(&node); err != nil {
		c.Error(errors.NewHTTPError(err, "failed to get node", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, node)
}

func (h *DefaultNodeHandler) HandlerGetHouseNodes(c *gin.Context) {
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to parse query(offset) to int", http.StatusBadRequest))
		return
	}

	houseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to parse param(id) to int", http.StatusBadRequest))
		return
	}

	nodes, count, err := h.NodeService.GetHouseNodes(houseID, offset)
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to get nodes", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Nodes": nodes,
		"Count": count,
	})
}

func (h *DefaultNodeHandler) HandlerGetNodes(c *gin.Context) {
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to parse query(offset) to int", http.StatusBadRequest))
		return
	}

	nodes, count, err := h.NodeService.GetNodes(offset)
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to get nodes", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Nodes": nodes,
		"Count": count,
	})
}
