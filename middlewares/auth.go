package middlewares

import (
	"github.com/buliqioqiolibusdo/demp-core/constants"
	"github.com/buliqioqiolibusdo/demp-core/controllers"
	"github.com/buliqioqiolibusdo/demp-core/errors"
	"github.com/buliqioqiolibusdo/demp-core/user"
	"github.com/gin-gonic/gin"
)

func AuthorizationMiddleware() gin.HandlerFunc {
	userSvc, _ := user.GetUserService()
	return func(c *gin.Context) {
		// token string
		tokenStr := c.GetHeader("Authorization")

		// validate token
		u, err := userSvc.CheckToken(tokenStr)
		if err != nil {
			// validation failed, return error response
			controllers.HandleErrorUnauthorized(c, errors.ErrorHttpUnauthorized)
			return
		}

		// set user in context
		c.Set(constants.ContextUser, u)

		// validation success
		c.Next()
	}
}
