package controllers

import (
	"github.com/buliqioqiolibusdo/demp-core/constants"
	"github.com/buliqioqiolibusdo/demp-core/interfaces"
	"github.com/gin-gonic/gin"
)

func GetUserFromContext(c *gin.Context) (u interfaces.User) {
	value, ok := c.Get(constants.ContextUser)
	if !ok {
		return nil
	}
	u, ok = value.(interfaces.User)
	if !ok {
		return nil
	}
	return u
}
