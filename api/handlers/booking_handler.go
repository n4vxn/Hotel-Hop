package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	api "github.com/n4vxn/Hotel-Hop/api/utils"
	"github.com/n4vxn/Hotel-Hop/db"
	"github.com/n4vxn/Hotel-Hop/types"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

func (h *BookingHandler) HandleGetBookings(c *gin.Context) {
	bookings, err := h.store.Booking.GetBookings()
	if err != nil {
		api.HandleError(c, api.ErrNotFound("bookings"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": bookings})
}

func (h *BookingHandler) HandleGetBooking(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, _ := strconv.Atoi(userIDStr)
	booking, err := h.store.Booking.GetBookingByID(userID)
	if err != nil {
		api.HandleError(c, api.ErrUnauthorized("bookings"))
		return
	}
	user, ok := c.Get("user")
	if !ok {
		api.HandleError(c, err)
		c.Abort()
		return
	}
	typedUser, ok := user.(*types.User)
	if !ok {
		api.HandleError(c, err)
		c.Abort()
		return
	}
	if typedUser.UserID != userID {
		api.HandleError(c, api.ErrUnauthorized("user"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": booking})
}

func (h *BookingHandler) HandleCancelBookings(c *gin.Context) {
	userIDstr := c.Param("id")
	userID, _ := strconv.Atoi(userIDstr)

	_, err := h.store.Booking.GetBookingByID(userID)
	if err != nil {
		api.HandleError(c, api.ErrUnauthorized("user"))
		return
	}

	user, ok := c.Get("user")
	if !ok {
		api.HandleError(c, err)
		c.Abort()
		return
	}

	typedUser, ok := user.(*types.User)
	if !ok {
		api.HandleError(c, err)
		c.Abort()
		return
	}

	if !typedUser.IsAdmin && typedUser.UserID != userID {
		api.HandleError(c, api.ErrUnauthorized("user"))
		return
	}

	err = h.store.Booking.DeleteBookingByID(userID)
	if err != nil {
		api.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Bookings cancelled succesfully"})
}
