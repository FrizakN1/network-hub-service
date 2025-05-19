package router

import (
	"backend/database"
	"backend/proto/userpb"
	"backend/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

type EventHandler interface {
	handlerGetEvents(c *gin.Context, from string)
}

type DefaultEventHandler struct {
	EventService database.EventService
	UserService  userpb.UserServiceClient
	Metadata     Metadata
}

func NewEventHandler(userClient *userpb.UserServiceClient) EventHandler {
	return &DefaultEventHandler{
		EventService: &database.DefaultEventService{},
		UserService:  *userClient,
		Metadata:     &DefaultMetadata{},
	}
}

func (h *DefaultEventHandler) handlerGetEvents(c *gin.Context, from string) {
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

	ctx := h.Metadata.setAuthorizationHeader(c)

	userResp, err := h.UserService.GetUsersByIds(ctx, &userpb.GetUsersByIdsRequest{Ids: userIds})
	if err != nil {
		utils.Logger.Println(err)
		handlerError(c, err, 500)
		return
	}

	usersMap := make(map[int32]*userpb.User)
	for _, user := range userResp.Users {
		usersMap[user.Id] = user
	}

	for i := range events {
		events[i].User = usersMap[events[i].UserId]
	}

	c.JSON(200, gin.H{
		"Events": events,
		"Count":  count,
	})
}
