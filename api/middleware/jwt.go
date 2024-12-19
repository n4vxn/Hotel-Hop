package middleware

import (
	"fmt"
	"os"
	"strconv"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	api "github.com/n4vxn/Hotel-Hop/api/utils"
	"github.com/n4vxn/Hotel-Hop/db"
	"github.com/n4vxn/Hotel-Hop/types"
)

func JWTkey() []byte {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file: %v\n", err)
	}

	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		fmt.Println("JWT_SECRET environment variable not set")
		return nil
	}

	return []byte(secretKey)
}
func JWTAuthMiddleware(userStore db.UserStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("X-Api-Key")
		if authHeader == "" {
			api.HandleError(c, api.ErrUnauthorizedHeader())
			c.Abort()
			return
		}

		claims, err := VerifyToken(authHeader)
		if err != nil {
			api.HandleError(c, api.ErrInvalidToken())
			c.Abort()
			return
		}

		email, ok := claims["email"].(string)
		if !ok {
			api.HandleError(c, api.ErrInvalidTokenClaims())
			c.Abort()
			return
		}

		userIDstr, ok := claims["userID"].(string)
		if !ok {
			api.HandleError(c, api.ErrInvalidToken())
			c.Abort()
			return
		}

		userID, _ := strconv.Atoi(userIDstr)
		user, err := userStore.GetUserByID(userID)
		if err != nil {
			api.HandleError(c, api.ErrUnauthorized("user"))
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Set("email", email)
		c.Next()
	}
}

func CreateTokenFromUser(user *types.User) string {
	now := time.Now()
	validUntil := now.Add(time.Hour * 1).Unix() // Token valid for 1 hour

	claims := jwt.MapClaims{
		"userID":  fmt.Sprintf("%d", user.UserID),
		"email":   user.Email,
		"expires": validUntil,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JWTkey())
	if err != nil {
		fmt.Println("failed to sign token with secret: %w", err)
	}
	return tokenString
}

func VerifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return JWTkey(), nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not parse token: %w", err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}
	return claims, nil
}
