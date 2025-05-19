package handlers

import (
	"backend/database"
	"backend/errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type AddressHandler interface {
	HandlerGetHouses(c *gin.Context)
	HandlerGetHouse(c *gin.Context)
	HandlerGetSuggestions(c *gin.Context)
}

type DefaultAddressHandler struct {
	AddressService database.AddressService
}

func NewAddressHandler() AddressHandler {
	return &DefaultAddressHandler{
		AddressService: &database.DefaultAddressService{},
	}
}

func (h *DefaultAddressHandler) HandlerGetHouses(c *gin.Context) {
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to parse query(offset) to int", http.StatusBadRequest))
		return
	}

	addresses, count, err := h.AddressService.GetHouses(offset)
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
	var address database.Address
	var err error

	address.House.ID, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to parse param(id) to int", http.StatusBadRequest))
		return
	}

	err = h.AddressService.GetHouse(&address)
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

	suggestions, count, err := h.AddressService.GetSuggestions(search, offset, limit)
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to get suggestions", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Addresses": suggestions,
		"Count":     count,
	})
}
