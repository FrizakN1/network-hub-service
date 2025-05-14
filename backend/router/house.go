package router

import (
	"backend/database"
	"backend/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

type AddressHandler interface {
	handlerGetHouses(c *gin.Context)
	handlerGetHouse(c *gin.Context)
	handlerGetSuggestions(c *gin.Context)
}

type DefaultAddressHandler struct {
	AddressService database.AddressService
}

func (ah *DefaultAddressHandler) handlerGetHouses(c *gin.Context) {
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	addresses, count, err := ah.AddressService.GetHouses(offset)
	if err != nil {
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, gin.H{
		"Addresses": addresses,
		"Count":     count,
	})
}

func (ah *DefaultAddressHandler) handlerGetHouse(c *gin.Context) {
	var address database.Address
	var err error

	address.House.ID, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	err = ah.AddressService.GetHouse(&address)
	if err != nil {
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, address)
}

func (ah *DefaultAddressHandler) handlerGetSuggestions(c *gin.Context) {
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		utils.Logger.Println(err)
	}
	search := c.DefaultQuery("search", "")

	suggestions, count, err := ah.AddressService.GetSuggestions(search, offset, limit)
	if err != nil {
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, gin.H{
		"Addresses": suggestions,
		"Count":     count,
	})
}

func NewAddressHandler() AddressHandler {
	return &DefaultAddressHandler{
		AddressService: &database.DefaultAddressService{},
	}
}
