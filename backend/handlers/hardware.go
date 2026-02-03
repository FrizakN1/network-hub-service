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

type HardwareHandler interface {
	HandlerGetHardwareByID(c *gin.Context)
	HandlerEditHardware(c *gin.Context)
	HandlerCreateHardware(c *gin.Context)
	HandlerGetSearchHardware(c *gin.Context)
	HandlerGetNodeHardware(c *gin.Context)
	HandlerGetHouseHardware(c *gin.Context)
	HandlerGetHardware(c *gin.Context)
	HandlerDeleteHardware(c *gin.Context)
	SendBatchHardware(ctx context.Context) error
	SendSingleHardware(ctx context.Context, hardwareID int) error
}

type DefaultHardwareHandler struct {
	Privilege      Privilege
	HardwareRepo   database.HardwareRepository
	EventRepo      database.EventRepository
	AddressService addresspb.AddressServiceClient
	Metadata       utils.Metadata
	SearchService  searchpb.SearchServiceClient
	kafka.HardwareProducer
	utils.Logger
}

func NewHardwareHandler(addressClient *addresspb.AddressServiceClient, searchClient *searchpb.SearchServiceClient, db *database.Database, logger *utils.Logger) HardwareHandler {
	return &DefaultHardwareHandler{
		Privilege: &DefaultPrivilege{},
		HardwareRepo: &database.DefaultHardwareRepository{
			Database: *db,
		},
		EventRepo: &database.DefaultEventRepository{
			Database: *db,
		},
		AddressService:   *addressClient,
		Metadata:         &utils.DefaultMetadata{},
		SearchService:    *searchClient,
		HardwareProducer: kafka.NewHardwareProducer(kafka.NewKafkaWriter("index-node")),
		Logger:           *logger,
	}
}

func (h *DefaultHardwareHandler) SendSingleHardware(ctx context.Context, hardwareID int) error {
	hd := &models.Hardware{ID: hardwareID}

	if err := h.HardwareRepo.GetHardwareByID(hd); err != nil {
		return err
	}

	res, err := h.AddressService.GetAddress(ctx, &addresspb.GetAddressRequest{HouseId: hd.Node.HouseId})
	if err != nil {
		return err
	}

	grpcHd := &searchpb.Hardware{
		Id:        int32(hd.ID),
		Type:      hd.Type.Value,
		NodeName:  hd.Node.Name,
		ModelName: hd.Switch.Name,
		IpAddress: hd.IpAddress.String,
		Address: &searchpb.Address{
			StreetName: res.Street.Name,
			StreetType: res.Street.Type.ShortName,
			HouseName:  res.House.Name,
			HouseType:  res.House.Type.ShortName,
		},
		IsDelete: hd.IsDelete,
	}

	if err = h.HardwareProducer.SendSingleHardware(ctx, grpcHd); err != nil {
		return err
	}

	return nil
}

func (h *DefaultHardwareHandler) SendBatchHardware(ctx context.Context) error {
	hardware, err := h.HardwareRepo.GetHardwareForIndex()
	if err != nil {
		return err
	}

	if len(hardware) == 0 {
		return nil
	}

	if err = h.getAddressesForHardware(ctx, hardware); err != nil {
		return err
	}

	var grpcHardware []*searchpb.Hardware

	for _, hd := range hardware {

		grpcHd := &searchpb.Hardware{
			Id:        int32(hd.ID),
			Type:      hd.Type.Value,
			NodeName:  hd.Node.Name,
			ModelName: hd.Switch.Name,
			IpAddress: hd.IpAddress.String,
			Address: &searchpb.Address{
				StreetName: hd.Node.Address.Street.Name,
				StreetType: hd.Node.Address.Street.Type.ShortName,
				HouseName:  hd.Node.Address.House.Name,
				HouseType:  hd.Node.Address.House.Type.ShortName,
			},
			IsDelete: hd.IsDelete,
		}

		grpcHardware = append(grpcHardware, grpcHd)
	}

	const batchSize = 1000

	for i := 0; i < len(grpcHardware); i += batchSize {
		end := i + batchSize
		if end > len(grpcHardware) {
			end = len(grpcHardware)
		}

		batch := grpcHardware[i:end]

		if err = h.HardwareProducer.SendBatchHardware(ctx, batch); err != nil {
			return err
		}
	}

	return nil
}

//func (h *DefaultHardwareHandler) HandlerIndexHardware(c *gin.Context) {
//	hardware, err := h.HardwareRepo.GetHardwareForIndex()
//	if err != nil {
//		c.Error(errors.NewHTTPError(err, "failed to get hardware for index", http.StatusInternalServerError))
//		return
//	}
//
//	ctx := h.Metadata.SetAuthorizationHeader(c)
//
//	if err = h.getAddressesForHardware(ctx, hardware); err != nil {
//		c.Error(errors.NewHTTPError(err, "failed to get addresses", http.StatusInternalServerError))
//		return
//	}
//
//	var grpcHardware []*searchpb.Hardware
//
//	for _, hd := range hardware {
//		grpcHd := &searchpb.Hardware{
//			Id:        int32(hd.ID),
//			Type:      hd.Type.Value,
//			NodeName:  hd.Node.Name,
//			ModelName: hd.Switch.Name,
//			IpAddress: hd.IpAddress.String,
//			Address: &searchpb.Address{
//				StreetName: hd.Node.Address.Street.Name,
//				StreetType: hd.Node.Address.Street.Type.ShortName,
//				HouseName:  hd.Node.Address.House.Name,
//				HouseType:  hd.Node.Address.House.Type.ShortName,
//			},
//			IsDelete: hd.IsDelete,
//		}
//
//		grpcHardware = append(grpcHardware, grpcHd)
//	}
//
//	_, err = h.SearchService.IndexHardware(ctx, &searchpb.IndexHardwareRequest{Hardware: grpcHardware})
//	if err != nil {
//		c.Error(errors.NewHTTPError(err, "failed to index hardware", http.StatusInternalServerError))
//		return
//	}
//
//	c.JSON(http.StatusOK, nil)
//}

func (h *DefaultHardwareHandler) HandlerDeleteHardware(c *gin.Context) {
	_, isAdmin, _ := h.Privilege.getPrivilege(c)

	if !isAdmin {
		c.Error(errors.NewHTTPError(nil, "forbidden", http.StatusForbidden))
		return
	}

	hardwareID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to parse param(id) to int", http.StatusBadRequest))
		return
	}

	if err = h.HardwareRepo.DeleteHardware(hardwareID); err != nil {
		c.Error(errors.NewHTTPError(err, "failed to delete hardware", http.StatusInternalServerError))
		return
	}

	go func() {
		if e := h.SendSingleHardware(context.Background(), hardwareID); e != nil {
			log.Printf("failed to send single hardware: %v\n", e)
			h.Logger.Println(e)
		}
	}()

	c.JSON(http.StatusOK, true)
}

func (h *DefaultHardwareHandler) HandlerGetHardwareByID(c *gin.Context) {
	var (
		err      error
		hardware models.Hardware
	)

	hardware.ID, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to parse param(id)", http.StatusBadRequest))
		return
	}

	if err = h.HardwareRepo.GetHardwareByID(&hardware); err != nil {
		c.Error(errors.NewHTTPError(err, "failed to get hardware", http.StatusBadRequest))
		return
	}

	ctx := h.Metadata.SetAuthorizationHeader(c)

	res, e := h.AddressService.GetAddress(ctx, &addresspb.GetAddressRequest{HouseId: hardware.Node.HouseId})
	if e != nil {
		c.Error(errors.NewHTTPError(e, "failed to get addresses", http.StatusInternalServerError))
		return
	}

	hardware.Node.Address = &addresspb.Address{
		Street: res.Street,
		House:  res.House,
	}

	c.JSON(http.StatusOK, hardware)
}

func (h *DefaultHardwareHandler) HandlerEditHardware(c *gin.Context) {
	session, _, isOperatorOrHigher := h.Privilege.getPrivilege(c)

	if !isOperatorOrHigher {
		c.Error(errors.NewHTTPError(nil, "forbidden", http.StatusForbidden))
		return
	}

	var hardware models.Hardware

	if err := c.BindJSON(&hardware); err != nil {
		c.Error(errors.NewHTTPError(err, "invalid json", http.StatusBadRequest))
		return
	}

	if !h.HardwareRepo.ValidateHardware(hardware) {
		c.Error(errors.NewHTTPError(nil, "invalid hardware data", http.StatusBadRequest))
		return
	}

	hardware.UpdatedAt = sql.NullInt64{Int64: time.Now().Unix(), Valid: true}

	if err := h.HardwareRepo.EditHardware(&hardware); err != nil {
		c.Error(errors.NewHTTPError(err, "failed to edit hardware", http.StatusInternalServerError))
		return
	}

	event := models.Event{
		HouseId:     hardware.Node.HouseId,
		Node:        &models.Node{ID: hardware.Node.ID},
		Hardware:    &models.Hardware{ID: hardware.ID},
		UserId:      session.User.Id,
		Description: fmt.Sprintf("Изменение оборудования: %s", hardware.Type.Value),
		CreatedAt:   time.Now().Unix(),
	}

	if err := h.EventRepo.CreateEvent(event); err != nil {
		c.Error(errors.NewHTTPError(err, "failed to create event", http.StatusInternalServerError))
		return
	}

	go func() {
		if e := h.SendSingleHardware(context.Background(), hardware.ID); e != nil {
			log.Printf("failed to send single hardware: %v\n", e)
			h.Logger.Println(e)
		}
	}()

	c.JSON(http.StatusOK, hardware)
}

func (h *DefaultHardwareHandler) HandlerCreateHardware(c *gin.Context) {
	session, _, isOperatorOrHigher := h.Privilege.getPrivilege(c)

	if !isOperatorOrHigher {
		c.Error(errors.NewHTTPError(nil, "forbidden", http.StatusForbidden))
		return
	}

	var hardware models.Hardware

	if err := c.BindJSON(&hardware); err != nil {
		c.Error(errors.NewHTTPError(err, "invalid json", http.StatusBadRequest))
		return
	}

	if !h.HardwareRepo.ValidateHardware(hardware) {
		c.Error(errors.NewHTTPError(nil, "invalid hardware data", http.StatusBadRequest))
		return
	}

	hardware.CreatedAt = time.Now().Unix()

	if err := h.HardwareRepo.CreateHardware(&hardware); err != nil {
		c.Error(errors.NewHTTPError(err, "failed to create hardware", http.StatusInternalServerError))
		return
	}

	event := models.Event{
		HouseId:     hardware.Node.HouseId,
		Node:        &models.Node{ID: hardware.Node.ID},
		Hardware:    nil,
		UserId:      session.User.Id,
		Description: fmt.Sprintf("Создание оборудования: %s", hardware.Type.Value),
		CreatedAt:   time.Now().Unix(),
	}

	if err := h.EventRepo.CreateEvent(event); err != nil {
		c.Error(errors.NewHTTPError(err, "failed to create event", http.StatusInternalServerError))
	}

	go func() {
		if e := h.SendSingleHardware(context.Background(), hardware.ID); e != nil {
			log.Printf("failed to send single hardware: %v\n", e)
			h.Logger.Println(e)
		}
	}()

	c.JSON(http.StatusOK, hardware)
}

func (h *DefaultHardwareHandler) HandlerGetSearchHardware(c *gin.Context) {
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to parse query(offset) to int", http.StatusBadRequest))
		return
	}
	search := c.Query("search")

	ctx := h.Metadata.SetAuthorizationHeader(c)

	res, err := h.SearchService.SearchHardware(ctx, &searchpb.SearchHardwareRequest{
		Search:       &searchpb.Search{Query: search, Offset: int32(offset), Limit: 20},
		SearchFilter: &searchpb.SearchHardwareFilter{UseIsDelete: false},
	})
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to search hardware", http.StatusInternalServerError))
		return
	}

	if res == nil || len(res.HardwareIDs) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"Hardware": []struct{}{},
			"Count":    0,
		})

		return
	}

	hardware, err := h.HardwareRepo.GetHardwareByIDs(res.HardwareIDs)
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to get hardware", http.StatusInternalServerError))
		return
	}

	if err = h.getAddressesForHardware(ctx, hardware); err != nil {
		c.Error(errors.NewHTTPError(err, "failed to get addresses", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Hardware": hardware,
		"Count":    res.Total,
	})
}

func (h *DefaultHardwareHandler) HandlerGetNodeHardware(c *gin.Context) {
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to parse query(offset) to int", http.StatusBadRequest))
		return
	}

	nodeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to parse param(id) to int", http.StatusBadRequest))
		return
	}

	hardware, count, err := h.HardwareRepo.GetHardware(offset, 0, nodeID)
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to get hardware", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Hardware": hardware,
		"Count":    count,
	})
}

func (h *DefaultHardwareHandler) HandlerGetHouseHardware(c *gin.Context) {
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

	hardware, count, err := h.HardwareRepo.GetHardware(offset, houseID, 0)
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to get hardware", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Hardware": hardware,
		"Count":    count,
	})
}

func (h *DefaultHardwareHandler) HandlerGetHardware(c *gin.Context) {
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to parse query(offset) to int", http.StatusBadRequest))
		return
	}

	hardware, count, err := h.HardwareRepo.GetHardware(offset, 0, 0)
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to get hardware", http.StatusBadRequest))
		return
	}

	houseIDSet := make(map[int32]struct{})
	addressMap := make(map[int32]*addresspb.Address)

	for _, hd := range hardware {
		houseIDSet[hd.Node.HouseId] = struct{}{}
	}

	var houseIDs []int32
	for houseID := range houseIDSet {
		houseIDs = append(houseIDs, houseID)
	}

	ctx := h.Metadata.SetAuthorizationHeader(c)

	res, err := h.AddressService.GetAddresses(ctx, &addresspb.GetAddressesRequest{HouseIDs: houseIDs})
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to get addresses", http.StatusInternalServerError))
		return
	}

	for _, address := range res.Addresses {
		addressMap[address.House.Id] = address
	}

	for i := range hardware {
		hardware[i].Node.Address = addressMap[hardware[i].Node.HouseId]
	}

	c.JSON(http.StatusOK, gin.H{
		"Hardware": hardware,
		"Count":    count,
	})
}

func (h *DefaultHardwareHandler) getAddressesForHardware(ctx context.Context, hardware []models.Hardware) error {
	houseIDSet := make(map[int32]struct{})
	addressMap := make(map[int32]*addresspb.Address)

	for _, hd := range hardware {
		houseIDSet[hd.Node.HouseId] = struct{}{}
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

	for i := range hardware {
		hardware[i].Node.Address = addressMap[hardware[i].Node.HouseId]
	}

	return nil
}
