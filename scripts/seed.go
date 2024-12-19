package scripts

import (
	"database/sql"
	"fmt"
	"github.com/n4vxn/Hotel-Hop/types"
)

func SeedHotelStore(db *sql.DB) error {
	hotels := []types.Hotel{
		{HotelName: "Sea View Hotel", Location: "Miami", Rating: 3},
		{HotelName: "Mountain Retreat", Location: "Aspen", Rating: 5},
	}

	for _, hotel := range hotels {
		_, err := db.Exec("INSERT INTO hotels (hotel_name, location, rating) VALUES ($1, $2, $3)", hotel.HotelName, hotel.Location, hotel.Rating)
		if err != nil {
			return fmt.Errorf("error inserting hotel: %v", err)
		}
	}
	return nil
}

func SeedRoomStore(db *sql.DB) error {
	// Sample rooms data
	rooms := []types.Room{
		{HotelID: 1, Size: "small", Price: 55.0},
		{HotelID: 1, Size: "large", Price: 85.0},
		{HotelID: 2, Size: "medium", Price: 120.0},
		{HotelID: 2, Size: "small", Price: 180.0},
	}

	for _, room := range rooms {
		_, err := db.Exec("INSERT INTO rooms (hotel_id, size, seaside, price) VALUES ($1, $2, $3, $4)", room.HotelID, room.Size, room.Seaside, room.Price)
		if err != nil {
			return fmt.Errorf("error inserting room: %v", err)
		}
	}

	return nil
}
