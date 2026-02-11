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
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
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
	HandlerGetNodesExcel(c *gin.Context)
}

type DefaultNodeHandler struct {
	Privilege      Privilege
	NodeRepo       database.NodeRepository
	ReportRepo     database.ReportRepository
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
		ReportRepo: &database.DefaultReportRepository{
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

func generateExcel(nodes []models.Node, reportData map[string]string) ([]byte, error) {
	sheetName := "Sheet1"

	f := excelize.NewFile()
	defer f.Close()

	_, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}

	if err = setExcelHeaders(f, sheetName); err != nil {
		return nil, err
	}

	if err = setExcelData(f, sheetName, nodes, reportData); err != nil {
		return nil, err
	}

	if err = setExcelStyle(f, sheetName); err != nil {
		return nil, err
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func setExcelStyle(f *excelize.File, sheetName string) error {
	if err := f.SetColWidth(sheetName, "A", "A", 40); err != nil {
		return err
	}

	if err := f.SetColWidth(sheetName, "B", "G", 35); err != nil {
		return err
	}

	shortCols := []string{"C", "E", "G"}

	for _, col := range shortCols {
		if err := f.SetColWidth(sheetName, col, col, 15); err != nil {
			return err
		}
	}

	styleCenter, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	if err != nil {
		return err
	}

	styleLeft, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center"},
	})
	if err != nil {
		return err
	}

	styleRight, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "right", Vertical: "center"},
	})
	if err != nil {
		return err
	}

	if err = f.SetCellStyle(sheetName, "A1", "G9", styleCenter); err != nil {
		return err
	}

	if err = f.SetCellStyle(sheetName, "A2", "A9", styleLeft); err != nil {
		return err
	}

	colLeft := []string{"B", "D", "F"}

	for _, col := range colLeft {
		if err = f.SetCellStyle(sheetName, fmt.Sprintf("%s5", col), fmt.Sprintf("%s6", col), styleLeft); err != nil {
			return err
		}

		if err = f.SetCellStyle(sheetName, fmt.Sprintf("%s7", col), fmt.Sprintf("%s7", col), styleRight); err != nil {
			return err
		}
	}

	return nil
}

func setExcelData(f *excelize.File, sheetName string, nodes []models.Node, reportData map[string]string) error {
	colNum := 2

	for i, node := range nodes {
		data := [][]interface{}{
			{fmt.Sprintf("Узел связи №%d", i+1)},
			{node.Placement.String},
			{node.Supply.String},
			{"модель", "мощность, кВт"},
			{reportData["HN_SWITCH"], reportData["HN_SWITCH_POWER"]},
			{reportData["OPTICAL_RECEIVER"], reportData["OPTICAL_RECEIVER_POWER"]},
			{"Итого мощность:", reportData["HN_TOTAL_CAPACITY"]},
			{reportData["VOLTAGE_LEVEL"]},
			{reportData["CATEGORY_RELIABILITY_POWER_SUPPLY"]},
		}

		if node.Type != nil && node.Type.Key != "HN" {
			data[4] = []interface{}{reportData["BN_SWITCH"], reportData["BN_SWITCH_POWER"]}
			data[6] = []interface{}{"Итого мощность:", reportData["BN_TOTAL_CAPACITY"]}
		}

		for j, row := range data {
			cell, _ := excelize.CoordinatesToCellName(colNum, j+1)
			cellNext, _ := excelize.CoordinatesToCellName(colNum+1, j+1)

			if err := f.SetCellValue(sheetName, cell, row[0]); err != nil {
				return err
			}

			if len(row) > 1 {
				if err := f.SetCellValue(sheetName, cellNext, row[1]); err != nil {
					return err
				}
			} else {
				if err := f.MergeCell(sheetName, cell, cellNext); err != nil {
					return err
				}
			}
		}

		colNum += 2
	}

	return nil
}

func setExcelHeaders(f *excelize.File, sheetName string) error {
	headers := []interface{}{
		"Средство связи",
		"Размещение узла",
		"Точки присоединения по эл.энергии",
		"Перечень и мощность оборудования",
		"",
		"",
		"",
		"Уровень напряжения",
		"Категория надежности электроснабжения",
	}

	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(1, i+1)

		if err := f.SetCellValue(sheetName, cell, header); err != nil {
			return err
		}
	}

	if err := f.MergeCell(sheetName, "A4", "A7"); err != nil {
		return err
	}

	return nil
}

func (h *DefaultNodeHandler) HandlerGetNodesExcel(c *gin.Context) {
	houseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to parse param(id) to int", http.StatusBadRequest))
		return
	}

	nodes, _, err := h.NodeRepo.GetNodes(0, true, houseID)
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to get nodes", http.StatusInternalServerError))
		return
	}

	reportData, err := h.ReportRepo.GetReportData()
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to get report data", http.StatusInternalServerError))
		return
	}

	dataMap, err := parseReportData(reportData)
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to parse report data", http.StatusInternalServerError))
	}

	excelData, err := generateExcel(nodes, dataMap)
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to generate Excel", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, base64.StdEncoding.EncodeToString(excelData))
}

func parseReportData(reportData []models.Report) (map[string]string, error) {
	reportDataMap := make(map[string]string)

	for _, report := range reportData {
		reportDataMap[report.Key] = report.Value
	}

	HNSwitchPower, err := strconv.ParseFloat(reportDataMap["HN_SWITCH_POWER"], 8)
	if err != nil {
		return nil, err
	}

	BNSwitchPower, err := strconv.ParseFloat(reportDataMap["BN_SWITCH_POWER"], 8)
	if err != nil {
		return nil, err
	}

	opticalReceiverPower, err := strconv.ParseFloat(reportDataMap["OPTICAL_RECEIVER_POWER"], 8)
	if err != nil {
		return nil, err
	}

	reportDataMap["HN_TOTAL_CAPACITY"] = fmt.Sprintf("%v", HNSwitchPower+opticalReceiverPower)
	reportDataMap["BN_TOTAL_CAPACITY"] = fmt.Sprintf("%v", BNSwitchPower+opticalReceiverPower)

	return reportDataMap, nil
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

	if node.Parent != nil {
		houseIDs := []int32{node.HouseId, node.Parent.HouseId}
		addressMap := make(map[int32]*addresspb.Address)

		res, e := h.AddressService.GetAddresses(ctx, &addresspb.GetAddressesRequest{HouseIDs: houseIDs})
		if e != nil {
			c.Error(errors.NewHTTPError(e, "failed to get addresses", http.StatusInternalServerError))
			return
		}

		for _, address := range res.Addresses {
			addressMap[address.House.Id] = address
		}

		node.Address = addressMap[node.HouseId]
		node.Parent.Address = addressMap[node.Parent.HouseId]
	} else {
		res, e := h.AddressService.GetAddress(ctx, &addresspb.GetAddressRequest{HouseId: node.HouseId})
		if e != nil {
			c.Error(errors.NewHTTPError(e, "failed to get addresses", http.StatusInternalServerError))
			return
		}

		node.Address = &addresspb.Address{
			Street: res.Street,
			House:  res.House,
		}
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
