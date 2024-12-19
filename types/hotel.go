package types

type RoomType int

const (
	SingleRoomType RoomType = iota
	DoubleRoomType
	SeaSideRoomType
	DeluxeRoomType
)

type Hotel struct {
	HotelID   int    `json:"hotel_id,omitempty"`
	HotelName string `json:"hotel_name"`
	Location  string `json:"location"`
	Rooms     []Room `json:"rooms"`
	Rating    int    `json:"rating"`
}

type CreateHotelParams struct {
	HotelID   int    `json:"hotel_id,omitempty"`
	HotelName string `json:"hotel_name"`
	Location  string `json:"location"`
	Rating    int    `json:"rating"`
}

type Room struct {
	RoomID    int     `json:"room_id,omitempty"`
	HotelID   int     `json:"hotel_id"`
	Size      string  `json:"size"`
	Seaside   bool    `json:"seaside"`
	Price     float64 `json:"price"`
	Available bool    `json:"available"`
}

func NewHotelFromParams(params Hotel) (*CreateHotelParams, error) {
	return &CreateHotelParams{
		HotelName: params.HotelName,
		Location:  params.Location,
		Rating:    params.Rating,
	}, nil
}
