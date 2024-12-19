package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/n4vxn/Hotel-Hop/api/middleware"
	api "github.com/n4vxn/Hotel-Hop/api/utils"
	"github.com/n4vxn/Hotel-Hop/db"
	"github.com/n4vxn/Hotel-Hop/types"
)

type AuthParams struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthHandler struct {
	userStore db.UserStore
}

type AuthRespose struct {
	User  *types.User `json:"user"`
	Token string      `json:"token"`
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

func (h *AuthHandler) HandleAuthenticate(c *gin.Context) {
	var Params AuthParams
	if err := c.ShouldBindJSON(&Params); err != nil {
		api.HandleError(c, api.ErrInvalidRequestData())
		return
	}
	user, err := h.userStore.GetUserByEmail(Params.Email)
	fmt.Println("user: ", user)
	if err != nil {
		api.HandleError(c, api.ErrNotFound("user"))
		return
	}
	err = types.IsValidPassword(user.EncryptedPassword, Params.Password)
	if err != nil {
		api.HandleError(c, api.ErrInvalidCredentials())
	}
	resp := AuthRespose{
		User:  user,
		Token: middleware.CreateTokenFromUser(user),
	}
	c.JSON(http.StatusOK, resp)
}
