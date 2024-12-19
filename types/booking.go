package types

import "time"

type Booking struct {
	BookingID  int       `json:"booking_id"`
	UserID     int       `json:"user_id"`
	RoomID     int       `json:"room_id"`
	NumPersons int       `json:"numpersons"`
	FromDate   time.Time `json:"fromDate,omitempty"`
	TillDate   time.Time `json:"tillDate,omitempty"`
}
