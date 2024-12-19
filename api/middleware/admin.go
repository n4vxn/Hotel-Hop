package middleware

import (
	"github.com/gin-gonic/gin"
	api "github.com/n4vxn/Hotel-Hop/api/utils"
	"github.com/n4vxn/Hotel-Hop/types"
)

func AdminAuth(c *gin.Context) {
	var err error
	user, exists := c.Get("user")
	if !exists {
		api.HandleError(c, api.ErrUnauthorized(""))
		c.Abort()
		return
	}

	typedUser, ok := user.(*types.User)
	if !ok {
		api.HandleError(c, err)
		c.Abort()
		return
	}

	if !typedUser.IsAdmin {
		api.HandleError(c, api.ErrIsAdmin())
		c.Abort()
		return
	}

	c.Next()
}
