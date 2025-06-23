package handlers

import (
	"backend/database"
	"backend/errors"
	"backend/models"
	"backend/proto/addresspb"
	"backend/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"os"
	"strconv"
)

type AddressHandler interface {
	HandlerGetHouse(c *gin.Context)
	HandlerSearchAddresses(c *gin.Context)
	HandlerGetHouses(c *gin.Context)
	HandlerSetHouseParams(c *gin.Context)
}

type DefaultAddressHandler struct {
	AddressService addresspb.AddressServiceClient
	Metadata       utils.Metadata
	AddressRepo    database.AddressRepository
}

func NewAddressHandler(addressClient *addresspb.AddressServiceClient, db *database.Database) AddressHandler {
	return &DefaultAddressHandler{
		AddressService: *addressClient,
		Metadata:       &utils.DefaultMetadata{},
		AddressRepo:    &database.DefaultAddressRepository{Database: *db},
	}
}

func InitAddressClient() *addresspb.AddressServiceClient {
	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%s", os.Getenv("ADDRESS_SERVICE_ADDRESS"), os.Getenv("ADDRESS_SERVICE_PORT")),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("could not connect to address service: %v", err)
	}

	addressClient := addresspb.NewAddressServiceClient(conn)

	return &addressClient
}

func (h *DefaultAddressHandler) HandlerSetHouseParams(c *gin.Context) {
	address := &models.AddressParams{}
	var err error

	if err = c.BindJSON(&address); err != nil {
		c.Error(errors.NewHTTPError(err, "invalid json", http.StatusBadRequest))
		return
	}

	address.HouseID, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to parse param(id) to int", http.StatusBadRequest))
		return
	}

	if err = h.AddressRepo.SetHouseParams(address); err != nil {
		c.Error(errors.NewHTTPError(err, "failed to set house params", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, address)
}

func (h *DefaultAddressHandler) HandlerGetHouse(c *gin.Context) {
	houseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to parse param(id) to int", http.StatusBadRequest))
		return
	}

	ctx := h.Metadata.SetAuthorizationHeader(c)

	res, err := h.AddressService.GetAddress(ctx, &addresspb.GetAddressRequest{HouseId: int32(houseID)})
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to get house", http.StatusInternalServerError))
		return
	}

	addressParams := &models.AddressParams{
		HouseID: houseID,
	}

	if err = h.AddressRepo.GetAddressParams(addressParams); err != nil {
		c.Error(errors.NewHTTPError(err, "failed to get address params", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Address": &addresspb.Address{
			Street: res.Street,
			House:  res.House,
		},
		"Params": addressParams,
	})
}

func (h *DefaultAddressHandler) HandlerGetHouses(c *gin.Context) {
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to parse query(offset) to int", http.StatusBadRequest))
		return
	}

	addressAmountsMap, err := h.AddressRepo.GetAddressesAmounts(nil, offset)
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to query houses", http.StatusInternalServerError))
		return
	}

	houseIDs := make([]int32, len(addressAmountsMap))
	count := 0

	for houseID, _ := range addressAmountsMap {
		houseIDs[count] = houseID
		count++
	}

	ctx := h.Metadata.SetAuthorizationHeader(c)

	res, err := h.AddressService.GetAddresses(ctx, &addresspb.GetAddressesRequest{HouseIDs: houseIDs})
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to get houses", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Addresses":      res.Addresses,
		"AddressAmounts": addressAmountsMap,
		"Count":          count,
	})
}

func (h *DefaultAddressHandler) HandlerSearchAddresses(c *gin.Context) {
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to parse query(offset) to int", http.StatusBadRequest))
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to parse query(limit) to int", http.StatusBadRequest))
		return
	}
	search := c.DefaultQuery("search", "")

	ctx := h.Metadata.SetAuthorizationHeader(c)

	res, err := h.AddressService.SearchAddresses(ctx, &addresspb.SearchAddressesRequest{
		Search: search,
		Offset: int32(offset),
		Limit:  int32(limit),
	})
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to search addresses", http.StatusInternalServerError))
		return
	}

	houseIDs := make([]int32, len(res.Addresses))
	for i, address := range res.Addresses {
		houseIDs[i] = address.House.Id
	}

	addressAmounts, err := h.AddressRepo.GetAddressesAmounts(houseIDs, 0)
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to get house params", http.StatusInternalServerError))
	}

	c.JSON(http.StatusOK, gin.H{
		"Addresses":      res.Addresses,
		"AddressAmounts": addressAmounts,
		"Count":          res.Total,
	})
}
