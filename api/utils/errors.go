package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Error struct {
	Code int    `json:"string"`
	Err  string `json:"err"`
}

func (e Error) Error() string {
	return e.Err
}

func NewError(code int, err string) Error {
	return Error{
		Code: code,
		Err:  err,
	}
}

func HandleError(c *gin.Context, err error) {
	switch e := err.(type) {
	case Error:
		c.JSON(e.Code, gin.H{"error": e.Err})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": e.Error()})
	}
}

//user
func ErrUnauthorized(resource string) Error {
	return NewError(http.StatusUnauthorized, "Unauthorized "+resource)
}

func ErrInvalidID() Error {
	return NewError(http.StatusBadRequest, "Invalid ID given")
}

func ErrInvalidRequestData() Error {
	return NewError(http.StatusBadRequest, "Invalid request data")
}

func ErrNotFound(resource string) Error {
	return NewError(http.StatusNotFound, resource+" not found")
}

func ErrCreation(resource string) Error {
	return NewError(http.StatusInternalServerError, "Failed to create "+resource)
}

func ErrUpdation(resource string) Error {
	return NewError(http.StatusInternalServerError, "Failed to update "+resource)
}

func ErrInvalidCredentials() Error {
	return NewError(http.StatusUnauthorized, "Invalid Credentials")
}

//hotel & room
func ErrRoomAvailabilityChecking() Error {
	return NewError(http.StatusInternalServerError, "Error checking room availability")
}

func ErrRoomNotAvailable() Error {
	return NewError(http.StatusInternalServerError, "Room is already booked for the requested dates")
}

func ErrBooking() Error {
	return NewError(http.StatusInternalServerError, "Failed to book room")
}

func ErrIsAdmin() Error {
	return NewError(http.StatusInternalServerError, "Forbidden: Admins only")
}

//jwt
func ErrUnauthorizedHeader() Error {
	return NewError(http.StatusUnauthorized, "Unauthorized request")
}

func ErrInvalidToken() Error {
	return NewError(http.StatusUnauthorized, "Invalid token")
}

func ErrInvalidTokenClaims() Error {
	return NewError(http.StatusUnauthorized, "Invalid token claims")
}

func ErrUserNotFound() Error {
	return NewError(http.StatusUnauthorized, "User not found")
}

