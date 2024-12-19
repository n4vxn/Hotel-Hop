package types

import (
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

var (
	minFirstNameLen = 2
	minLastNameLen  = 2
	minPassLen      = 7
)

type CreateUserParams struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	IsAdmin   bool   `json:"is_admin"`
}

type UpdateUserParams struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type User struct {
	UserID            int    `json:"user_id,omitempty"`
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	Email             string `json:"email"`
	EncryptedPassword string `json:"-"`
	IsAdmin           bool   `json:"is_admin"`
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"error"`
}

func IsValidPassword(encpw, pw string) error {
	return bcrypt.CompareHashAndPassword([]byte(encpw), []byte(pw))
}

func (params CreateUserParams) Validate() []ValidationError {
	var errors []ValidationError
	if len(params.FirstName) < minFirstNameLen {
		errors = append(errors, ValidationError{"first_name", fmt.Sprintf("first name should be at least %d characters", minFirstNameLen)})
	}
	if len(params.LastName) < minLastNameLen {
		errors = append(errors, ValidationError{"last_name", fmt.Sprintf("last name should be atleast %d characters", minLastNameLen)})
	}
	if !isEmailValid(params.Email) {
		errors = append(errors, ValidationError{"email", "email is invalid"})
	}
	return errors
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(e)
}
func NewUserFromParams(params CreateUserParams) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &User{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: string(encpw),
		IsAdmin:           params.IsAdmin,
	}, nil
}
