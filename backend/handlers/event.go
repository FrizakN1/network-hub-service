package handlers

import (
	"backend/database"
	"backend/errors"
	"backend/proto/addresspb"
	"backend/proto/userpb"
	"backend/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type EventHandler interface {
	HandlerGetEvents(c *gin.Context, from string)
}

type DefaultEventHandler struct {
	EventRepo      database.EventRepository
	UserService    userpb.UserServiceClient
	AddressService addresspb.AddressServiceClient
	Metadata       utils.Metadata
}

func NewEventHandler(userClient *userpb.UserServiceClient, addressClient *addresspb.AddressServiceClient, db *database.Database) EventHandler {
	return &DefaultEventHandler{
		EventRepo: &database.DefaultEventRepository{
			Database: *db,
		},
		UserService:    *userClient,
		AddressService: *addressClient,
		Metadata:       &utils.DefaultMetadata{},
	}
}

func (h *DefaultEventHandler) HandlerGetEvents(c *gin.Context, from string) {
	id := 0
	var err error

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to parse query(offset) to int", http.StatusBadRequest))
		return
	}

	if from != "" {
		from = fmt.Sprintf("%s_%s", from, strings.ToUpper(c.Param("type")))

		id, err = strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Error(errors.NewHTTPError(err, "failed to parse param(id) to int", http.StatusBadRequest))
			return
		}
	}

	events, count, err := h.EventRepo.GetEvents(offset, from, id)
	if err != nil {
		c.Error(errors.NewHTTPError(err, "failed to get events", http.StatusInternalServerError))
		return
	}

	userIDSet := make(map[int32]struct{})
	houseIDSet := make(map[int32]struct{})
	usersMap := make(map[int32]*userpb.User)
	addressMap := make(map[int32]*addresspb.Address)

	for _, event := range events {
		userIDSet[event.UserId] = struct{}{}
		houseIDSet[event.HouseId] = struct{}{}
	}

	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	ctx := h.Metadata.SetAuthorizationHeader(c)

	if len(userIDSet) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()

			var userIDs []int32
			for userID := range userIDSet {
				userIDs = append(userIDs, userID)
			}

			userRes, e := h.UserService.GetUsersByIds(ctx, &userpb.GetUsersByIdsRequest{Ids: userIDs})
			if e != nil {
				errChan <- errors.NewHTTPError(e, "failed to get users", http.StatusInternalServerError)
				return
			}

			for _, user := range userRes.Users {
				usersMap[user.Id] = user
			}
		}()
	}

	if len(houseIDSet) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()

			var houseIDs []int32
			for houseID := range houseIDSet {
				houseIDs = append(houseIDs, houseID)
			}

			addressRes, e := h.AddressService.GetAddresses(ctx, &addresspb.GetAddressesRequest{HouseIDs: houseIDs})
			if e != nil {
				errChan <- errors.NewHTTPError(e, "failed to get addresses", http.StatusInternalServerError)
				return
			}

			for _, address := range addressRes.Addresses {
				addressMap[address.House.Id] = address
			}
		}()
	}

	wg.Wait()
	close(errChan)

	for e := range errChan {
		if e != nil {
			c.Error(e)
			return
		}
	}

	for i := range events {
		events[i].User = usersMap[events[i].UserId]
		events[i].Address = addressMap[events[i].HouseId]
	}

	c.JSON(http.StatusOK, gin.H{
		"Events": events,
		"Count":  count,
	})
}
