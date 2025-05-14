package router

import (
	"backend/database"
	"backend/proto/userpb"
	"backend/utils"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

type NodeHandler interface {
	handlerGetSearchNodes(c *gin.Context)
	handlerEditNode(c *gin.Context)
	handlerCreateNode(c *gin.Context)
	handlerGetNode(c *gin.Context)
	handlerGetHouseNodes(c *gin.Context)
	handlerGetNodes(c *gin.Context)
	handlerDeleteNode(c *gin.Context)
}

type DefaultNodeHandler struct {
	NodeService database.NodeService
}

func (h *DefaultHandler) handlerDeleteNode(c *gin.Context) {
	session, ok := c.Get("session")
	if !ok {
		err := errors.New("сессия не найдена")
		utils.Logger.Println(err)
		handlerError(c, err, 401)
		return
	}

	if session.(userpb.Session).User.Role.Value != "admin" {
		c.JSON(403, nil)
		return
	}

	nodeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	if err = h.NodeService.DeleteNode(nodeID); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, true)
}

func (h *DefaultHandler) handlerGetSearchNodes(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	search := c.Query("search")

	nodes, count, err := h.NodeService.GetSearchNodes(search, offset)
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, gin.H{
		"Nodes": nodes,
		"Count": count,
	})
}

func (h *DefaultHandler) handlerEditNode(c *gin.Context) {
	session, ok := c.Get("session")
	if !ok {
		err := errors.New("сессия не найдена")
		utils.Logger.Println(err)
		handlerError(c, err, 401)
		return
	}

	if session.(userpb.Session).User.Role.Value != "admin" && session.(userpb.Session).User.Role.Value != "operator" {
		c.JSON(403, nil)
		return
	}

	var node database.Node

	if err := c.BindJSON(&node); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	if !h.NodeService.ValidateNode(node) {
		c.JSON(400, nil)
		return
	}

	node.UpdatedAt = sql.NullInt64{
		Int64: time.Now().Unix(),
		Valid: true,
	}

	if err := h.NodeService.EditNode(&node); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	event := database.Event{
		Address:     database.Address{House: database.AddressElement{ID: node.Address.House.ID}},
		Node:        &database.Node{ID: node.ID},
		Hardware:    nil,
		User:        userpb.User{Id: session.(userpb.Session).User.Id},
		Description: fmt.Sprintf("Изменение узла: %s", node.Name),
		CreatedAt:   time.Now().Unix(),
	}

	if err := h.EventService.CreateEvent(event); err != nil {
		utils.Logger.Println(err)
	}

	c.JSON(200, node)
}

func (h *DefaultHandler) handlerCreateNode(c *gin.Context) {
	session, ok := c.Get("session")
	if !ok {
		err := errors.New("сессия не найдена")
		utils.Logger.Println(err)
		handlerError(c, err, 401)
		return
	}

	if session.(userpb.Session).User.Role.Value != "admin" && session.(userpb.Session).User.Role.Value != "operator" {
		c.JSON(403, nil)
		return
	}

	var node database.Node

	if err := c.BindJSON(&node); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	if !h.NodeService.ValidateNode(node) {
		c.JSON(400, nil)
		return
	}

	node.CreatedAt = time.Now().Unix()

	if err := h.NodeService.CreateNode(&node); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	event := database.Event{
		Address:     database.Address{House: database.AddressElement{ID: node.Address.House.ID}},
		Node:        nil,
		Hardware:    nil,
		User:        userpb.User{Id: session.(userpb.Session).User.Id},
		Description: fmt.Sprintf("Создание нового узла: %s", node.Name),
		CreatedAt:   time.Now().Unix(),
	}

	if err := h.EventService.CreateEvent(event); err != nil {
		utils.Logger.Println(err)
	}

	c.JSON(200, node)
}

func (h *DefaultHandler) handlerGetNode(c *gin.Context) {
	var (
		err  error
		node database.Node
	)

	node.ID, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	if err = h.NodeService.GetNode(&node); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, node)
}

func (h *DefaultHandler) handlerGetHouseNodes(c *gin.Context) {
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	houseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	nodes, count, err := h.NodeService.GetHouseNodes(houseID, offset)
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, gin.H{
		"Nodes": nodes,
		"Count": count,
	})
}

func (h *DefaultHandler) handlerGetNodes(c *gin.Context) {
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	nodes, count, err := h.NodeService.GetNodes(offset)
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, gin.H{
		"Nodes": nodes,
		"Count": count,
	})
}
