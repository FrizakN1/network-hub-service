package router

import (
	"backend/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

//func (h *DefaultHandler) handlerGetEventsFrom(c *gin.Context, from string) {
//	id, err := strconv.Atoi(c.Param("id"))
//	if err != nil {
//		utils.Logger.Println(err)
//		handlerError(c, err, 400)
//		return
//	}
//
//	events, count, err := h.EventService.GetEvents(from+"_"+strings.ToUpper(c.Param("type")), id)
//	if err != nil {
//		utils.Logger.Println(err)
//		handlerError(c, err, 400)
//		return
//	}
//
//	c.JSON(200, gin.H{
//		"Events": events,
//		"Count":  count,
//	})
//}

func (h *DefaultHandler) handlerGetEvents(c *gin.Context, from string) {
	id := 0
	var err error

	if from != "" {
		from = fmt.Sprintf("%s_%s", from, strings.ToUpper(c.Param("type")))

		id, err = strconv.Atoi(c.Param("id"))
		if err != nil {
			utils.Logger.Println(err)
			handlerError(c, err, 400)
			return
		}
	}

	events, count, err := h.EventService.GetEvents(from, id)
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 400)
		return
	}

	userIdSet := make(map[int32]struct{})
	for _, event := range events {
		userIdSet[event.UserId] = struct{}{}
	}

	var userIds []int32
	for userID := range userIdSet {
		userIds = append(userIds, userID)
	}

	usersMap, err := h.handlerGetUsersByIds(c, userIds)
	if err != nil {
		utils.Logger.Println(err)
		return
	}

	for i := range events {
		events[i].User = usersMap[events[i].UserId]
	}

	c.JSON(200, gin.H{
		"Events": events,
		"Count":  count,
	})
}
