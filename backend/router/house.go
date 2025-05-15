package router

import (
	"backend/database"
	"backend/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

func (h *DefaultHandler) handlerGetHouses(c *gin.Context) {
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	addresses, count, err := h.AddressService.GetHouses(offset)
	if err != nil {
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, gin.H{
		"Addresses": addresses,
		"Count":     count,
	})
}

func (h *DefaultHandler) handlerGetHouse(c *gin.Context) {
	var address database.Address
	var err error

	address.House.ID, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	err = h.AddressService.GetHouse(&address)
	if err != nil {
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, address)
}

func (h *DefaultHandler) handlerGetSuggestions(c *gin.Context) {
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

	suggestions, count, err := h.AddressService.GetSuggestions(search, offset, limit)
	if err != nil {
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, gin.H{
		"Addresses": suggestions,
		"Count":     count,
	})
}
