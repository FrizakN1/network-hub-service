package handlers

import (
	"backend/database"
	"backend/errors"
	"backend/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type AddressHandler interface {
	HandlerGetHouses(c *gin.Context)
	HandlerGetHouse(c *gin.Context)
	HandlerGetSuggestions(c *gin.Context)
	HandlerSetHouseParams(c *gin.Context)
}

type DefaultAddressHandler struct {
	AddressRepo database.AddressRepository
}

func NewAddressHandler(db *database.Database) AddressHandler {
	addressRepo := &database.DefaultAddressRepository{
		Database: *db,
	}

	addressRepo.LoadAddressElementTypeMap()

	return &DefaultAddressHandler{
		AddressRepo: addressRepo,
	}
}

func (h *DefaultAddressHandler) HandlerSetHouseParams(c *gin.Context) {
	var address models.Address
	var err error

	if err = c.BindJSON(&address); err != nil {
		c.Error(errors.NewHTTPError(err, "invalid json", http.StatusBadRequest))
		return
	}

	address.House.ID, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to parse param(id) to int", http.StatusBadRequest))
		return
	}

	if err = h.AddressRepo.SetHouseParams(&address); err != nil {
		c.Error(errors.NewHTTPError(err, "failed to set house params", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, address)
}

func (h *DefaultAddressHandler) HandlerGetHouses(c *gin.Context) {
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to parse query(offset) to int", http.StatusBadRequest))
		return
	}

	addresses, count, err := h.AddressRepo.GetHouses(offset)
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to get houses", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Addresses": addresses,
		"Count":     count,
	})
}

func (h *DefaultAddressHandler) HandlerGetHouse(c *gin.Context) {
	var address models.Address
	var err error

	address.House.ID, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to parse param(id) to int", http.StatusBadRequest))
		return
	}

	err = h.AddressRepo.GetHouse(&address)
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to get house", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, address)
}

func (h *DefaultAddressHandler) HandlerGetSuggestions(c *gin.Context) {
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

	suggestions, count, err := h.AddressRepo.GetSuggestions(search, offset, limit)
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to get suggestions", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Addresses": suggestions,
		"Count":     count,
	})
}
