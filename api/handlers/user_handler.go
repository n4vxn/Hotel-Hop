package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	api "github.com/n4vxn/Hotel-Hop/api/utils"
	"github.com/n4vxn/Hotel-Hop/db"
	"github.com/n4vxn/Hotel-Hop/types"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

// Handler to get a user by ID
func (h *UserHandler) HandlerGetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	user, err := h.userStore.GetUserByID(id)
	if err != nil {
		api.HandleError(c, api.ErrNotFound("user"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": user})
}

// Handler to get all users
func (h *UserHandler) HandlerGetUsers(c *gin.Context) {
	users, err := h.userStore.GetUsers()
	if err != nil {
		api.HandleError(c, api.ErrNotFound("users"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": users})
}

// Handler to create a new user
func (h *UserHandler) HandlePostUsers(c *gin.Context) {
	var params types.CreateUserParams
	err := c.ShouldBindJSON(&params)
	if err != nil {
		api.HandleError(c, api.ErrInvalidRequestData())
		return
	}

	if errors := params.Validate(); len(errors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": errors})
		return
	}

	user, err := types.NewUserFromParams(params)
	if err != nil {
		api.HandleError(c, api.ErrCreation("user"))
		return
	}

	err = h.userStore.CreateUser(user)
	if err != nil {
		api.HandleError(c, api.ErrNotFound("user in database"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User Created!"})
}

// Handler to delete a user by ID
func (h *UserHandler) HandleDeleteUsers(c *gin.Context) {
	id := c.Param("id")
	err := h.userStore.DeleteUser(id)
	if err != nil {
		if err.Error() == fmt.Sprintf("no user found with id %s", id) {
			api.HandleError(c, api.ErrNotFound("user"))
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}

// Handler to update user information
func (h *UserHandler) HandleUpdateUsers(c *gin.Context) {
	var user types.UpdateUserParams
	id := c.Param("id")
	err := c.ShouldBindJSON(&user)
	if err != nil {
		api.HandleError(c, api.ErrInvalidRequestData())
		return
	}

	err = h.userStore.UpdateUser(id, user)
	if err != nil {
		api.HandleError(c, api.ErrUpdation("user"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user updated successfully"})
}
