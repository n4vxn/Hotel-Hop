package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	api "github.com/n4vxn/Hotel-Hop/api/utils"
	"github.com/n4vxn/Hotel-Hop/db"
	"github.com/n4vxn/Hotel-Hop/types"
)

type HotelHandler struct {
	store *db.Store
}

func NewHotelHandler(store *db.Store) *HotelHandler {
	return &HotelHandler{
		store: store,
	}
}

// Handler to create a hotel
func (h *HotelHandler) HandleCreateHotel(c *gin.Context) {
	var params types.Hotel
	err := c.ShouldBindJSON(&params)
	if err != nil {
		api.HandleError(c, api.ErrInvalidRequestData())
		return
	}

	hotel, err := types.NewHotelFromParams(params)
	if err != nil {
		api.HandleError(c, api.ErrCreation("hotels"))
		return
	}

	err = h.store.Hotel.InsertHotel(hotel)
	if err != nil {
		api.HandleError(c, api.ErrCreation("hotel"))
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Hotel created successfully", "hotel": hotel})
}

// Handler to get rooms by hotel ID
func (h *HotelHandler) HandleGetRooms(c *gin.Context) {
	id := c.Param("id")
	rooms, err := h.store.Room.GetRoomsByHotelID(id)
	if err != nil {
		api.HandleError(c, api.ErrNotFound("rooms"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"rooms": rooms})
}

// Handler to get all hotels
func (h *HotelHandler) HandleGetHotels(c *gin.Context) {
	hotels, err := h.store.Hotel.GetHotels()
	if err != nil {
		api.HandleError(c, api.ErrNotFound("hotels"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": hotels})
}

// Handler to get a hotel by ID
func (h *HotelHandler) HandleGetHotelByID(c *gin.Context) {
	id := c.Param("id")
	hotel, err := h.store.Hotel.GetHotelsByID(id)
	if err != nil {
		api.HandleError(c, api.ErrNotFound("hotel"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": hotel})
}
