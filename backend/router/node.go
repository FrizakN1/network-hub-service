package router

import (
	"backend/database"
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

func (nh *DefaultNodeHandler) handlerDeleteNode(c *gin.Context) {
	sessionHash, ok := c.Get("sessionHash")
	if !ok {
		err := errors.New("сессия не найдена")
		utils.Logger.Println(err)
		handlerError(c, err, 401)
		return
	}

	session := database.GetSession(sessionHash.(string))

	if session.User.Role.Value != "admin" {
		c.JSON(403, nil)
		return
	}

	nodeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	if err = nh.NodeService.DeleteNode(nodeID); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, true)
}

func (nh *DefaultNodeHandler) handlerGetSearchNodes(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	search := c.Query("search")

	nodes, count, err := nh.NodeService.GetSearchNodes(search, offset)
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

func (nh *DefaultNodeHandler) handlerEditNode(c *gin.Context) {
	sessionHash, ok := c.Get("sessionHash")
	if !ok {
		err := errors.New("сессия не найдена")
		utils.Logger.Println(err)
		handlerError(c, err, 401)
		return
	}

	session := database.GetSession(sessionHash.(string))

	if session.User.Role.Value != "admin" && session.User.Role.Value != "operator" {
		c.JSON(403, nil)
		return
	}

	var node database.Node

	if err := c.BindJSON(&node); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	if !nh.NodeService.ValidateNode(node) {
		c.JSON(400, nil)
		return
	}

	node.UpdatedAt = sql.NullInt64{
		Int64: time.Now().Unix(),
		Valid: true,
	}

	if err := nh.NodeService.EditNode(&node); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	event := database.Event{
		Address:     database.Address{House: database.AddressElement{ID: node.Address.House.ID}},
		Node:        &database.Node{ID: node.ID},
		Hardware:    nil,
		User:        database.User{ID: session.User.ID},
		Description: fmt.Sprintf("Изменение узла: %s", node.Name),
		CreatedAt:   time.Now().Unix(),
	}

	if err := event.CreateEvent(); err != nil {
		utils.Logger.Println(err)
	}

	c.JSON(200, node)
}

func (nh *DefaultNodeHandler) handlerCreateNode(c *gin.Context) {
	sessionHash, ok := c.Get("sessionHash")
	if !ok {
		err := errors.New("сессия не найдена")
		utils.Logger.Println(err)
		handlerError(c, err, 401)
		return
	}

	session := database.GetSession(sessionHash.(string))

	if session.User.Role.Value != "admin" && session.User.Role.Value != "operator" {
		c.JSON(403, nil)
		return
	}

	var node database.Node

	if err := c.BindJSON(&node); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	if !nh.NodeService.ValidateNode(node) {
		c.JSON(400, nil)
		return
	}

	node.CreatedAt = time.Now().Unix()

	if err := nh.NodeService.CreateNode(&node); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	event := database.Event{
		Address:     database.Address{House: database.AddressElement{ID: node.Address.House.ID}},
		Node:        nil,
		Hardware:    nil,
		User:        database.User{ID: session.User.ID},
		Description: fmt.Sprintf("Создание нового узла: %s", node.Name),
		CreatedAt:   time.Now().Unix(),
	}

	if err := event.CreateEvent(); err != nil {
		utils.Logger.Println(err)
	}

	c.JSON(200, node)
}

func (nh *DefaultNodeHandler) handlerGetNode(c *gin.Context) {
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

	if err = nh.NodeService.GetNode(&node); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, node)
}

func (nh *DefaultNodeHandler) handlerGetHouseNodes(c *gin.Context) {
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

	nodes, count, err := nh.NodeService.GetHouseNodes(houseID, offset)
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

func NewNodeHandler() NodeHandler {
	return &DefaultNodeHandler{
		NodeService: &database.DefaultNodeService{},
	}
}

func (nh *DefaultNodeHandler) handlerGetNodes(c *gin.Context) {
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	nodes, count, err := nh.NodeService.GetNodes(offset)
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
