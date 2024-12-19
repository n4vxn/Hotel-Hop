package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/n4vxn/Hotel-Hop/types"
)

type UserStore interface {
	GetUserByID(int) (*types.User, error)
	GetUserByEmail(string) (*types.User, error)
	GetUsers() ([]*types.User, error)
	CreateUser(*types.User) error
	DeleteUser(string) error
	UpdateUser(string, types.UpdateUserParams) error
}

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{
		db: db,
	}
}

func (s *PostgresUserStore) InitUsersTable() error {
	return s.createUsersTable()
}

func (s *PostgresUserStore) createUsersTable() error {
	query := `CREATE TABLE IF NOT EXISTS users(
	user_id SERIAL PRIMARY KEY,
	first_name TEXT NOT NULL,
	last_name TEXT NOT NULL,
	email TEXT NOT NULL UNIQUE,
	encrypted_password TEXT NOT NULL,
	is_admin BOOLEAN DEFAULT FALSE)`

	_, err := s.db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating users table: %v", err)
	}
	return nil
}

func (s *PostgresUserStore) GetUserByID(id int) (*types.User, error) {
	query := `SELECT user_id, first_name, last_name, email, is_admin FROM users WHERE user_id = $1`
	row := s.db.QueryRow(query, id)

	var user types.User
	err := row.Scan(&user.UserID, &user.FirstName, &user.LastName, &user.Email, &user.IsAdmin)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *PostgresUserStore) GetUserByEmail(email string) (*types.User, error) {
	query := `SELECT user_id, first_name, last_name, email, encrypted_password, is_admin FROM users WHERE email = $1`
	row := s.db.QueryRow(query, email)

	var user types.User
	err := row.Scan(&user.UserID, &user.FirstName, &user.LastName, &user.Email, &user.EncryptedPassword, &user.IsAdmin)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (s *PostgresUserStore) GetUsers() ([]*types.User, error) {
	var users []*types.User
	query := `SELECT user_id, first_name, last_name, email FROM users`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		user := &types.User{}
		if err = rows.Scan(&user.UserID, &user.FirstName, &user.LastName, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (s *PostgresUserStore) CreateUser(user *types.User) error {
	fmt.Printf("%v\n", user)
	query := `INSERT INTO users (first_name, last_name, email, encrypted_password, is_admin) VALUES ($1, $2, $3, $4, $5) RETURNING user_id`
	err := s.db.QueryRow(query, user.FirstName, user.LastName, user.Email, user.EncryptedPassword, user.IsAdmin).Scan(&user.UserID)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresUserStore) DeleteUser(id string) error {
	query := `DELETE FROM users WHERE user_id = $1`
	result, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no user found with id %s", id)
	}

	return nil
}

func (s *PostgresUserStore) UpdateUser(id string, user types.UpdateUserParams) error {
	fmt.Printf("%v\n", user.FirstName)

	query := `UPDATE users SET first_name = COALESCE(NULLIF($1, ''), first_name), last_name = COALESCE(NULLIF($2, ''), last_name) WHERE user_id = $3`
	_, err := s.db.Exec(query, user.FirstName, user.LastName, id)
	if err != nil {
		return err
	}
	return nil
}
