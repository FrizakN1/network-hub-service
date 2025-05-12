package router

import (
	"backend/database"
	"backend/utils"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

func handlerGetEventsFrom(c *gin.Context, from string) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	events, count, err := database.GetEvents(from+"_"+strings.ToUpper(c.Param("type")), id)
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, gin.H{
		"Events": events,
		"Count":  count,
	})
}

func handlerGetEvents(c *gin.Context) {
	events, count, err := database.GetEvents("", 0)
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	c.JSON(200, gin.H{
		"Events": events,
		"Count":  count,
	})
}
