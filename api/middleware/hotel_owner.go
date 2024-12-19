package middleware

// import (
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// )

// func HotelOwnerAuth(c *gin.Context) {
// 	user, exists := c.Get("user")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
// 		c.Abort()
// 		return
// 	}

// 	typedUser, ok := user.(*types.)
// 	if !ok {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
// 		c.Abort()
// 		return
// 	}
// }