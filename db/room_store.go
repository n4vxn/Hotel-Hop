package db

import (
	"database/sql"
	"fmt"
	"github.com/n4vxn/Hotel-Hop/types"
)

type RoomStore interface {
	// InsertRoom(*types.Room) error
	GetRoomsByHotelID(id string) ([]*types.Room, error)
}

type PostgresRoomStore struct {
	db *sql.DB
}

func NewPostgresRoomStore(db *sql.DB) *PostgresRoomStore {
	return &PostgresRoomStore{
		db: db,
	}
}

func (s *PostgresRoomStore) InitRoomsTable() error {
	return s.createRoomsTable()
}

func (s *PostgresRoomStore) createRoomsTable() error {
	query := `CREATE TABLE IF NOT EXISTS rooms (
    room_id SERIAL PRIMARY KEY,
    hotel_id INT REFERENCES hotels(hotel_id) ON DELETE CASCADE,
    size TEXT NOT NULL,
	seaside boolean,
    price DECIMAL(10, 2) NOT NULL
	)`

	_, err := s.db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating rooms table: %v", err)
	}
	return nil
}

func (s *PostgresRoomStore) GetRoomsByHotelID(id string) ([]*types.Room, error) {
	var rooms []*types.Room
	query := `SELECT room_id, hotel_id, size, seaside, price FROM rooms WHERE hotel_id = $1`
	rows, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		room := types.Room{}
		if err := rows.Scan(&room.RoomID, &room.HotelID, &room.Size, &room.Seaside, &room.Price); err != nil {
			return nil, err
		}
		rooms = append(rooms, &room)
	}
	return rooms, err

}
