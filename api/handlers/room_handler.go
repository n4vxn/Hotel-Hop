package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	api "github.com/n4vxn/Hotel-Hop/api/utils"
	"github.com/n4vxn/Hotel-Hop/db"
	"github.com/n4vxn/Hotel-Hop/types"
)

type BookRoomParams struct {
	FromDate   time.Time `json:"fromDate"`
	TillDate   time.Time `json:"tillDate"`
	NumPersons int       `json:"numPersons"`
}

// Validate that booking dates are valid
func (p BookRoomParams) ValidateDate() error {
	now := time.Now()
	if now.After(p.FromDate) || now.After(p.TillDate) {
		return fmt.Errorf("invalid Date")
	}
	return nil
}

type RoomHandler struct {
	store db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: *store,
	}
}

// Handler to book a room
func (h *RoomHandler) HandleBookRoom(c *gin.Context) {
	var err error
	var params BookRoomParams
	if err := c.ShouldBindJSON(&params); err != nil {
		return
	}
	if err := params.ValidateDate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	roomIDstr := c.Param("id")
	roomID, _ := strconv.Atoi(roomIDstr)
	user, exists := c.Get("user")
	if !exists {
		api.HandleError(c, err)
		return
	}
	typedUser, ok := user.(*types.User)
	if !ok {
		api.HandleError(c, err)
		return
	}

	// Check if the room is available
	booking := types.Booking{
		UserID:     typedUser.UserID,
		RoomID:     roomID,
		FromDate:   params.FromDate,
		TillDate:   params.TillDate,
		NumPersons: params.NumPersons,
	}

	roomAvailable, err := h.store.Booking.RoomAvailability(roomID, booking.FromDate, booking.TillDate)
	if err != nil {
		api.HandleError(c, api.ErrRoomAvailabilityChecking())
		return
	}
	if !roomAvailable {
		api.HandleError(c, api.ErrRoomNotAvailable())
		return
	}

	// Insert the booking into the database
	insertedBooking, err := h.store.Booking.InsertBooking(&booking)
	if err != nil {
		api.HandleError(c, api.ErrBooking())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Room booked successfully", "booking": insertedBooking})
}
