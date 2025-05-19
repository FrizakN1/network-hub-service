package handlers

import (
	"backend/database"
	"backend/errors"
	"backend/proto/userpb"
	"backend/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

type EventHandler interface {
	HandlerGetEvents(c *gin.Context, from string)
}

type DefaultEventHandler struct {
	EventService database.EventService
	UserService  userpb.UserServiceClient
	Metadata     utils.Metadata
}

func NewEventHandler(userClient *userpb.UserServiceClient) EventHandler {
	return &DefaultEventHandler{
		EventService: &database.DefaultEventService{},
		UserService:  *userClient,
		Metadata:     &utils.DefaultMetadata{},
	}
}

func (h *DefaultEventHandler) HandlerGetEvents(c *gin.Context, from string) {
	id := 0
	var err error

	if from != "" {
		from = fmt.Sprintf("%s_%s", from, strings.ToUpper(c.Param("type")))

		id, err = strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Error(errors.NewHTTPError(err, "failed to parse param(id) to int", http.StatusBadRequest))
			return
		}
	}

	events, count, err := h.EventService.GetEvents(from, id)
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to get events", http.StatusInternalServerError))
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

	ctx := h.Metadata.SetAuthorizationHeader(c)

	userResp, err := h.UserService.GetUsersByIds(ctx, &userpb.GetUsersByIdsRequest{Ids: userIds})
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to get users", http.StatusInternalServerError))
		return
	}

	usersMap := make(map[int32]*userpb.User)
	for _, user := range userResp.Users {
		usersMap[user.Id] = user
	}

	for i := range events {
		events[i].User = usersMap[events[i].UserId]
	}

	c.JSON(http.StatusOK, gin.H{
		"Events": events,
		"Count":  count,
	})
}
