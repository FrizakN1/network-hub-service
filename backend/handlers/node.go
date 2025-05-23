package handlers

import (
	"backend/database"
	"backend/errors"
	"backend/models"
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
	Privilege Privilege
	NodeRepo  database.NodeRepository
	EventRepo database.EventRepository
}

func NewNodeHandler(db *database.Database) NodeHandler {
	return &DefaultNodeHandler{
		Privilege: &DefaultPrivilege{},
		NodeRepo: &database.DefaultNodeRepository{
			Database: *db,
		},
		EventRepo: &database.DefaultEventRepository{
			Database: *db,
		},
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

	if err = h.NodeRepo.DeleteNode(nodeID); err != nil {
		c.Error(errors.NewHTTPError(err, "failed to delete node", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, true)
}

func (h *DefaultNodeHandler) HandlerGetSearchNodes(c *gin.Context) {
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to parse query(offset) to int", http.StatusBadRequest))
		return
	}

	onlyActive, err := strconv.ParseBool(c.DefaultQuery("only_active", "false"))
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to parse query(only_active) to bool", http.StatusBadRequest))
		return
	}

	search := c.Query("search")

	nodes, count, err := h.NodeRepo.GetSearchNodes(search, offset, onlyActive)
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to get search nodes", http.StatusInternalServerError))
		return
	}
	//ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	//defer cancel()
	//
	//var wg sync.WaitGroup
	//var nodes []models.Node
	//var count int
	//var errChan = make(chan *errors.HTTPError)
	//
	//wg.Add(2)
	//
	//go func() {
	//	defer wg.Done()
	//
	//	var e error
	//	nodes, e = h.NodeRepo.GetSearchNodes(search, offset)
	//	if e != nil {
	//		errChan <- errors.NewHTTPError(e, "failed to get search nodes", http.StatusInternalServerError)
	//	}
	//	errChan <- nil
	//}()
	//
	//go func() {
	//	defer wg.Done()
	//
	//	var e error
	//	count, e = h.Counter.CountRecords(ctx, "SEARCH_NODES", []interface{}{search})
	//	if e != nil {
	//		errChan <- errors.NewHTTPError(e, "failed to get count search nodes", http.StatusInternalServerError)
	//	}
	//	errChan <- nil
	//}()
	//
	//go func() {
	//	wg.Wait()
	//	close(errChan)
	//}()
	//
	//for e := range errChan {
	//	if e != nil {
	//		cancel()
	//		c.Error(e)
	//		return
	//	}
	//}

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

	var node models.Node

	if err := c.BindJSON(&node); err != nil {
		c.Error(errors.NewHTTPError(err, "invalid json", http.StatusBadRequest))
		return
	}

	if !h.NodeRepo.ValidateNode(node) {
		c.Error(errors.NewHTTPError(nil, "invalid node data", http.StatusBadRequest))
		return
	}

	node.UpdatedAt = sql.NullInt64{
		Int64: time.Now().Unix(),
		Valid: true,
	}

	if err := h.NodeRepo.EditNode(&node); err != nil {
		c.Error(errors.NewHTTPError(err, "failed to edit node", http.StatusInternalServerError))
		return
	}

	event := models.Event{
		Address:     models.Address{House: models.AddressElement{ID: node.Address.House.ID}},
		Node:        &models.Node{ID: node.ID},
		Hardware:    nil,
		UserId:      session.User.Id,
		Description: fmt.Sprintf("Изменение узла: %s", node.Name),
		CreatedAt:   time.Now().Unix(),
	}

	if err := h.EventRepo.CreateEvent(event); err != nil {
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

	var node models.Node

	if err := c.BindJSON(&node); err != nil {
		c.Error(errors.NewHTTPError(err, "invalid json", http.StatusBadRequest))
		return
	}

	if !h.NodeRepo.ValidateNode(node) {
		c.Error(errors.NewHTTPError(nil, "invalid node data", http.StatusBadRequest))
		return
	}

	node.CreatedAt = time.Now().Unix()

	if err := h.NodeRepo.CreateNode(&node); err != nil {
		c.Error(errors.NewHTTPError(err, "failed to create node", http.StatusInternalServerError))
		return
	}

	event := models.Event{
		Address:     models.Address{House: models.AddressElement{ID: node.Address.House.ID}},
		Node:        nil,
		Hardware:    nil,
		UserId:      session.User.Id,
		Description: fmt.Sprintf("Создание нового узла: %s", node.Name),
		CreatedAt:   time.Now().Unix(),
	}

	if err := h.EventRepo.CreateEvent(event); err != nil {
		c.Error(errors.NewHTTPError(err, "failed to create event", http.StatusInternalServerError))
	}

	c.JSON(http.StatusOK, node)
}

func (h *DefaultNodeHandler) HandlerGetNode(c *gin.Context) {
	var (
		err  error
		node models.Node
	)

	node.ID, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to parse param(id) to int", http.StatusBadRequest))
		return
	}

	if err = h.NodeRepo.GetNode(&node); err != nil {
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
	nodes, count, err := h.NodeRepo.GetNodes(offset, false, houseID)
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

	onlyActive, err := strconv.ParseBool(c.DefaultQuery("only_active", "false"))
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to parse query(only_active) to bool", http.StatusBadRequest))
		return
	}

	nodes, count, err := h.NodeRepo.GetNodes(offset, onlyActive, 0)
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to get nodes", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Nodes": nodes,
		"Count": count,
	})
}
