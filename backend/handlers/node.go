package handlers

import (
	"backend/database"
	"backend/errors"
	"backend/kafka"
	"backend/models"
	"backend/proto/addresspb"
	"backend/proto/searchpb"
	"backend/utils"
	"context"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
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
	SendBatchNodes(ctx context.Context) error
	SendSingleNode(ctx context.Context, nodeID int) error
}

type DefaultNodeHandler struct {
	Privilege      Privilege
	NodeRepo       database.NodeRepository
	EventRepo      database.EventRepository
	AddressService addresspb.AddressServiceClient
	Metadata       utils.Metadata
	SearchService  searchpb.SearchServiceClient
	kafka.NodeProducer
	utils.Logger
}

func NewNodeHandler(addressClient *addresspb.AddressServiceClient, searchClient *searchpb.SearchServiceClient, db *database.Database, logger *utils.Logger) NodeHandler {
	return &DefaultNodeHandler{
		Privilege: &DefaultPrivilege{},
		NodeRepo: &database.DefaultNodeRepository{
			Database: *db,
		},
		EventRepo: &database.DefaultEventRepository{
			Database: *db,
		},
		AddressService: *addressClient,
		Metadata:       &utils.DefaultMetadata{},
		SearchService:  *searchClient,
		NodeProducer:   kafka.NewNodeProducer(kafka.NewKafkaWriter("index-node")),
		Logger:         *logger,
	}
}

func (h *DefaultNodeHandler) SendSingleNode(ctx context.Context, nodeID int) error {
	node := &models.Node{ID: nodeID}

	if err := h.NodeRepo.GetNode(node); err != nil {
		return err
	}

	res, err := h.AddressService.GetAddress(ctx, &addresspb.GetAddressRequest{HouseId: node.HouseId})
	if err != nil {
		return err
	}

	nodeType := "Активный"

	if node.IsPassive {
		nodeType = "Пассивный"
	}

	grpcNode := &searchpb.Node{
		Id:    int32(node.ID),
		Name:  node.Name,
		Zone:  node.Zone.String,
		Owner: node.Owner.Value,
		Address: &searchpb.Address{
			StreetName: res.Street.Name,
			StreetType: res.Street.Type.ShortName,
			HouseName:  res.House.Name,
			HouseType:  res.House.Type.ShortName,
		},
		Type:      nodeType,
		IsDelete:  node.IsDelete,
		IsPassive: node.IsPassive,
	}

	if err = h.NodeProducer.SendSingleNode(ctx, grpcNode); err != nil {
		return err
	}

	return nil
}

func (h *DefaultNodeHandler) SendBatchNodes(ctx context.Context) error {
	nodes, err := h.NodeRepo.GetNodesForIndex()
	if err != nil {
		return err
	}

	if len(nodes) == 0 {
		return nil
	}

	if err = h.getAddressesForNodes(ctx, nodes); err != nil {
		return err
	}

	var grpcNodes []*searchpb.Node

	for _, node := range nodes {
		nodeType := "Активный"

		if node.IsPassive {
			nodeType = "Пассивный"
		}

		grpcNode := &searchpb.Node{
			Id:    int32(node.ID),
			Name:  node.Name,
			Zone:  node.Zone.String,
			Owner: node.Owner.Value,
			Address: &searchpb.Address{
				StreetName: node.Address.Street.Name,
				StreetType: node.Address.Street.Type.ShortName,
				HouseName:  node.Address.House.Name,
				HouseType:  node.Address.House.Type.ShortName,
			},
			Type:      nodeType,
			IsDelete:  node.IsDelete,
			IsPassive: node.IsPassive,
		}

		grpcNodes = append(grpcNodes, grpcNode)
	}

	const batchSize = 1000

	for i := 0; i < len(grpcNodes); i += batchSize {
		end := i + batchSize
		if end > len(grpcNodes) {
			end = len(grpcNodes)
		}

		batch := grpcNodes[i:end]

		if err = h.NodeProducer.SendBatchNodes(ctx, batch); err != nil {
			return err
		}
	}

	return nil
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

	go func() {
		if e := h.SendSingleNode(context.Background(), nodeID); e != nil {
			log.Printf("failed to send single node: %v\n", e)
			h.Logger.Println(e)
		}
	}()

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

	ctx := h.Metadata.SetAuthorizationHeader(c)

	res, err := h.SearchService.SearchNodes(ctx, &searchpb.SearchNodesRequest{
		Search:       &searchpb.Search{Query: search, Offset: int32(offset), Limit: 20},
		SearchFilter: &searchpb.SearchNodeFilter{UseIsDelete: true, UseIsPassive: onlyActive, IsPassive: !onlyActive, IsDelete: false},
	})
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to search nodes", http.StatusInternalServerError))
		return
	}

	if res == nil || len(res.NodesIDs) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"Nodes": []struct{}{},
			"Count": 0,
		})

		return
	}

	nodes, err := h.NodeRepo.GetNodesByIDs(res.NodesIDs)
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to get nodes", http.StatusInternalServerError))
		return
	}

	if err = h.getAddressesForNodes(ctx, nodes); err != nil {
		c.Error(errors.NewHTTPError(err, "failed to get addresses", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Nodes": nodes,
		"Count": res.Total,
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
		HouseId:     node.HouseId,
		Node:        &models.Node{ID: node.ID},
		Hardware:    nil,
		UserId:      session.User.Id,
		Description: fmt.Sprintf("Изменение узла: %s", node.Name),
		CreatedAt:   time.Now().Unix(),
	}

	if err := h.EventRepo.CreateEvent(event); err != nil {
		c.Error(errors.NewHTTPError(err, "failed to delete node", http.StatusInternalServerError))
	}

	go func() {
		if e := h.SendSingleNode(context.Background(), node.ID); e != nil {
			log.Printf("failed to send single node: %v\n", e)
			h.Logger.Println(e)
		}
	}()

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
		HouseId:     node.HouseId,
		Node:        nil,
		Hardware:    nil,
		UserId:      session.User.Id,
		Description: fmt.Sprintf("Создание нового узла: %s", node.Name),
		CreatedAt:   time.Now().Unix(),
	}

	if err := h.EventRepo.CreateEvent(event); err != nil {
		c.Error(errors.NewHTTPError(err, "failed to create event", http.StatusInternalServerError))
	}

	go func() {
		if e := h.SendSingleNode(context.Background(), node.ID); e != nil {
			log.Printf("failed to send single node: %v\n", e)
			h.Logger.Println(e)
		}
	}()

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

	ctx := h.Metadata.SetAuthorizationHeader(c)

	res, e := h.AddressService.GetAddress(ctx, &addresspb.GetAddressRequest{HouseId: node.HouseId})
	if e != nil {
		c.Error(errors.NewHTTPError(e, "failed to get addresses", http.StatusInternalServerError))
		return
	}

	node.Address = &addresspb.Address{
		Street: res.Street,
		House:  res.House,
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

	ctx := h.Metadata.SetAuthorizationHeader(c)

	if err = h.getAddressesForNodes(ctx, nodes); err != nil {
		c.Error(errors.NewHTTPError(err, "failed to get addresses", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Nodes": nodes,
		"Count": count,
	})
}

func (h *DefaultNodeHandler) getAddressesForNodes(ctx context.Context, nodes []models.Node) error {
	houseIDSet := make(map[int32]struct{})
	addressMap := make(map[int32]*addresspb.Address)

	for _, node := range nodes {
		houseIDSet[node.HouseId] = struct{}{}
	}

	var houseIDs []int32
	for houseID := range houseIDSet {
		houseIDs = append(houseIDs, houseID)
	}

	res, err := h.AddressService.GetAddresses(ctx, &addresspb.GetAddressesRequest{HouseIDs: houseIDs})
	if err != nil {
		return err
	}

	for _, address := range res.Addresses {
		addressMap[address.House.Id] = address
	}

	for i := range nodes {
		nodes[i].Address = addressMap[nodes[i].HouseId]
	}

	return nil
}
