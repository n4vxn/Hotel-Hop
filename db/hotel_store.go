package db

import (
	"database/sql"
	"fmt"
	"github.com/n4vxn/Hotel-Hop/types"
)

type HotelStore interface {
	InsertHotel(hotel *types.CreateHotelParams) error
	GetHotels() ([]*types.Hotel, error)
	GetHotelsByID(id string) (*types.Hotel, error)
}

type PostgresHotelStore struct {
	db *sql.DB
}

func NewPostgresHotelStore(db *sql.DB) *PostgresHotelStore {
	return &PostgresHotelStore{
		db: db,
	}
}

func (s *PostgresHotelStore) InitHotelsTable() error {
	return s.createHotelsTable()
}

func (s *PostgresHotelStore) createHotelsTable() error {
	query := `CREATE TABLE IF NOT EXISTS hotels(
	hotel_id SERIAL PRIMARY KEY,
	hotel_name TEXT NOT NULL,
	location TEXT NOT NULL,
	rating int
	)`

	_, err := s.db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating hotels table: %v", err)
	}
	return nil
}

func (s *PostgresHotelStore) InsertHotel(hotel *types.CreateHotelParams) error {
	query := `INSERT INTO hotels(hotel_name, location, rating) VALUES($1, $2, $3) RETURNING hotel_id `
	err := s.db.QueryRow(query, hotel.HotelName, hotel.Location, hotel.Rating).Scan(&hotel.HotelID)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresHotelStore) GetHotels() ([]*types.Hotel, error) {
	query := `
		SELECT 
			h.hotel_id, h.hotel_name, h.location, h.rating,
			r.room_id, r.size, r.seaside, r.price
		FROM 
			hotels h
		LEFT JOIN 
			rooms r ON h.hotel_id = r.hotel_id`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error fetching hotels with rooms: %v", err)
	}
	defer rows.Close()

	var hotels []*types.Hotel
	hotelMap := make(map[int]*types.Hotel)

	for rows.Next() {
		var hotelID, roomID, rating int
		var hotelName, location, size string
		var price float64
		var seaside bool

		if err := rows.Scan(
			&hotelID, &hotelName, &location, &rating,
			&roomID, &size, &seaside, &price,
		); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		if hotel, exists := hotelMap[hotelID]; exists {
			if roomID > 0 {
				room := types.Room{
					RoomID:  roomID,
					Size:    size,
					Seaside: seaside,
					Price:   price,
				}
				hotel.Rooms = append(hotel.Rooms, room)
			}
		} else {
			hotel := &types.Hotel{
				HotelID:   hotelID,
				HotelName: hotelName,
				Location:  location,
				Rating:    rating,
				Rooms:     []types.Room{},
			}

			if roomID > 0 {
				room := types.Room{
					RoomID:  roomID,
					Size:    size,
					Seaside: seaside,
					Price:   price,
				}
				hotel.Rooms = append(hotel.Rooms, room)
			}

			hotelMap[hotelID] = hotel
		}
	}

	for _, hotel := range hotelMap {
		hotels = append(hotels, hotel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error with rows: %v", err)
	}

	return hotels, nil
}

func (s *PostgresHotelStore) GetHotelsByID(id string) (*types.Hotel, error) {
	query := `SELECT hotel_id, hotel_name, location, room, rating FROM hotels WHERE hotel_id = $1`

	var hotel types.Hotel

	err := s.db.QueryRow(query, id).Scan(&hotel.HotelID, &hotel.HotelName, &hotel.Location, &hotel.Rooms, &hotel.Rating)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &hotel, nil
}
