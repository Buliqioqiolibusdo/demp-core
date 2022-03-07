package controllers

import (
	"net/http"

	"github.com/buliqioqiolibusdo/demp-core/constants"
	"github.com/buliqioqiolibusdo/demp-core/errors"
	"github.com/buliqioqiolibusdo/demp-core/interfaces"
	"github.com/buliqioqiolibusdo/demp-core/models/models"
	"github.com/buliqioqiolibusdo/demp-core/user"
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
)

var LoginController ActionController

func getLoginActions() []Action {
	loginCtx := newLoginContext()
	return []Action{
		{Method: http.MethodPost, Path: "/login", HandlerFunc: loginCtx.login},
		{Method: http.MethodPost, Path: "/logout", HandlerFunc: loginCtx.logout},
	}
}

type loginContext struct {
	userSvc interfaces.UserService
}

func (ctx *loginContext) login(c *gin.Context) {
	var u models.User
	if err := c.ShouldBindJSON(&u); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	token, loggedInUser, err := ctx.userSvc.Login(&interfaces.UserLoginOptions{
		Username: u.Username,
		Password: u.Password,
	})
	if err != nil {
		HandleErrorUnauthorized(c, errors.ErrorUserUnauthorized)
		return
	}
	c.Set(constants.UserContextKey, loggedInUser)
	HandleSuccessWithData(c, token)
}

func (ctx *loginContext) logout(c *gin.Context) {
	c.Set(constants.UserContextKey, nil)
	HandleSuccess(c)
}

func newLoginContext() *loginContext {
	// context
	ctx := &loginContext{}

	// dependency injection
	c := dig.New()
	if err := c.Provide(user.ProvideGetUserService()); err != nil {
		panic(err)
	}
	if err := c.Invoke(func(
		userSvc interfaces.UserService,
	) {
		ctx.userSvc = userSvc
	}); err != nil {
		panic(err)
	}

	return ctx
}
