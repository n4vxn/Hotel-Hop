package db

import (
	"database/sql"
	"fmt"
	"time"
	"github.com/n4vxn/Hotel-Hop/types"
)

type BookingStore interface {
	InsertBooking(*types.Booking) (*types.Booking, error)
	RoomAvailability(roomID int, fromDate time.Time, tillDate time.Time) (bool, error)
	GetBookings() ([]*types.Booking, error)
	GetBookingByID(user_id int) ([]*types.Booking, error)
	DeleteBookingByID(user_id int) error
}

type PostgresBookingStore struct {
	db *sql.DB
}

func NewPostgresBookingStore(db *sql.DB) *PostgresBookingStore {
	return &PostgresBookingStore{
		db: db,
	}
}

func (s *PostgresBookingStore) InitBookingsTable() error {
	return s.createBookingsTable()
}

func (s *PostgresBookingStore) createBookingsTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS bookings (
		booking_id SERIAL PRIMARY KEY,
		user_id INT NOT NULL REFERENCES users(user_id),
		room_id INT NOT NULL REFERENCES rooms(room_id),
		from_date TIMESTAMP NOT NULL,
		till_date TIMESTAMP NOT NULL,
		num_persons INT NOT NULL
	)`

	_, err := s.db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating bookings table: %v", err)
	}
	return nil
}

func (s *PostgresBookingStore) InsertBooking(booking *types.Booking) (*types.Booking, error) {
	query := `
		INSERT INTO bookings (user_id, room_id, from_date, till_date, num_persons)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING booking_id`
	fmt.Println(booking)
	row := s.db.QueryRow(query, booking.UserID, booking.RoomID, booking.FromDate, booking.TillDate, booking.NumPersons)
	err := row.Scan(&booking.BookingID)
	if err != nil {
		return nil, fmt.Errorf("error inserting booking: %v", err)
	}

	fmt.Println(booking)
	return booking, nil
}

func (s *PostgresBookingStore) RoomAvailability(roomID int, fromDate time.Time, tillDate time.Time) (bool, error) {
	query := `
        SELECT EXISTS(
            SELECT 1
            FROM bookings
            WHERE room_id = $1
              AND (
                from_date <= $3 AND till_date >= $2
              )
        );
    `
	var exists bool
	err := s.db.QueryRow(query, roomID, fromDate, tillDate).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	return !exists, nil
}

func (s *PostgresBookingStore) GetBookings() ([]*types.Booking, error) {
	var bookings []*types.Booking
	query := `
		SELECT booking_id, user_id, room_id, from_date, till_date, num_persons FROM bookings`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		booking := &types.Booking{}
		err := rows.Scan(&booking.BookingID, &booking.UserID, &booking.RoomID, &booking.FromDate, &booking.TillDate, &booking.NumPersons)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, booking)
	}
	return bookings, nil
}

func (s *PostgresBookingStore) GetBookingByID(user_id int) ([]*types.Booking, error) {
	var bookings []*types.Booking
	query := `
		SELECT booking_id, user_id, room_id, from_date, till_date, num_persons FROM bookings WHERE user_id = $1`
	rows, err := s.db.Query(query, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		booking := &types.Booking{}
		err := rows.Scan(&booking.BookingID, &booking.UserID, &booking.RoomID, &booking.FromDate, &booking.TillDate, &booking.NumPersons)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, fmt.Errorf("no booking found with id %d", user_id)
			}
			return nil, fmt.Errorf("error retrieving booking: %v", err)
		}
		bookings = append(bookings, booking)
	}
	return bookings, nil
}

func (s *PostgresBookingStore) DeleteBookingByID(user_id int) error {
	query := `
		DELETE FROM bookings WHERE user_id = $1`
	row, _ := s.db.Exec(query, user_id)
	if row != nil {
		return nil
	}
	return nil
}
