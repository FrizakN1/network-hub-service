package router

import (
	"backend/database"
	"backend/utils"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"time"
)

func handlerGetSearchNodes(c *gin.Context) {
	request := struct {
		Text   string
		Offset int
	}{}

	if err := c.BindJSON(&request); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	nodes, count, err := database.GetSearchNodes(request.Text, request.Offset)
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

func handlerGetNodeImages(c *gin.Context) {
	nodeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	files, err := database.GetNodeFiles(nodeID, true)
	if err != nil {
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, files)
}

func handlerGetNodeFiles(c *gin.Context) {
	nodeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	files, err := database.GetNodeFiles(nodeID, false)
	if err != nil {
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, files)
}

func handlerEditNode(c *gin.Context) {
	var node database.Node

	if err := c.BindJSON(&node); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	if !node.ValidateNode() {
		c.JSON(400, nil)
		return
	}

	node.UpdatedAt = sql.NullInt64{
		Int64: time.Now().Unix(),
		Valid: true,
	}

	if err := node.EditNode(); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, node)
}

func handlerCreateNode(c *gin.Context) {
	var node database.Node

	if err := c.BindJSON(&node); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	if !node.ValidateNode() {
		c.JSON(400, nil)
		return
	}

	node.CreatedAt = time.Now().Unix()

	if err := node.CreateNode(); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, node)
}

func handlerEditReferenceRecord(c *gin.Context) {
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

	var owner database.Enum

	if err := c.BindJSON(&owner); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	if len(owner.Name) == 0 {
		c.JSON(400, nil)
		return
	}

	if err := owner.EditReferenceRecord(strings.ToUpper(c.Param("reference"))); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, owner)
}

func handlerCreateReferenceRecord(c *gin.Context) {
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

	var referenceRecord database.Enum
	reference := strings.ToUpper(c.Param("reference"))

	if err := c.BindJSON(&referenceRecord); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	switch reference {
	case "NODE_TYPE":
	case "OWNER":
		if len(referenceRecord.Name) == 0 {
			c.JSON(400, nil)
			return
		}
		break
	case "HARDWARE_TYPE":
	case "OPERATION_MODE":
		if len(referenceRecord.Value) == 0 || len(referenceRecord.TranslateValue) == 0 {
			c.JSON(400, nil)
			return
		}
		break
	}

	referenceRecord.CreatedAt = time.Now().Unix()

	if err := referenceRecord.CreateReferenceRecord(reference); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, referenceRecord)
}

func handlerGetNodeTypes(c *gin.Context) {
	nodeTypes, err := database.GetNodeEnums("NODE_TYPES")
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, nodeTypes)
}

func handlerGetOwners(c *gin.Context) {
	owners, err := database.GetNodeEnums("OWNERS")
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, owners)
}

func handlerGetNode(c *gin.Context) {
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

	if err = node.GetNode(); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, node)
}

func handlerGetHouseNodes(c *gin.Context) {
	request := struct {
		Offset int
	}{}

	if err := c.BindJSON(&request); err != nil {
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

	nodes, count, err := database.GetHouseNodes(houseID, request.Offset)
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

func handlerGetNodes(c *gin.Context) {
	request := struct {
		Offset int
	}{}

	if err := c.BindJSON(&request); err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	nodes, count, err := database.GetNodes(request.Offset)
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
