package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-core/models/service"
	"github.com/crawlab-team/crawlab-core/user"
	"github.com/crawlab-team/crawlab-core/utils"
	"github.com/crawlab-team/go-trace"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/dig"
)

var UserController *userController

func getUserActions() []Action {
	userCtx := newUserContext()
	return []Action{
		{
			Method:      http.MethodPost,
			Path:        "/:id/change-password",
			HandlerFunc: userCtx.changePassword,
		},
		{
			Method:      http.MethodGet,
			Path:        "/me",
			HandlerFunc: userCtx.getMe,
		},
	}
}

func (ctx *userContext) getMe(c *gin.Context) {
	u, err := ctx._getMe(c)
	if err != nil {
		HandleErrorUnauthorized(c, errors.ErrorUserUnauthorized)
		return
	}
	HandleSuccessWithData(c, u)
}

func (ctx *userContext) _getMe(c *gin.Context) (u interfaces.User, err error) {
	res, ok := c.Get(constants.ContextUser)
	if !ok { // by zhizhong
		return nil, trace.TraceError(errors.ErrorUserNotExistsInContext)
	}
	u, ok = res.(interfaces.User)
	if !ok {
		return nil, trace.TraceError(errors.ErrorUserInvalidType)
	}
	// u.GetRole()
	return u, nil
}

type userController struct {
	ListActionControllerDelegate
	d   ListActionControllerDelegate
	ctx *userContext
}

// Create user
func (ctr *userController) Put(c *gin.Context) {
	var u models.User
	// var currentuser interfaces.User
	if err := c.ShouldBindJSON(&u); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// // Check root user by zhizhong
	tokenStr := c.GetHeader("Authorization")
	userSvc, _ := user.GetUserService()
	u_info, _ := userSvc.CheckToken(tokenStr)
	if u_info.GetRole() != "root" {
		HandleErrorUnauthorized(c, errors.ErrorHttpUnauthorized)
		return
	}

	if err := ctr.ctx.userSvc.Create(&interfaces.UserCreateOptions{
		Username: u.Username,
		Password: u.Password,
		Email:    u.Email,
		Role:     u.Role,
	}); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccess(c)
}

func (ctr *userController) PostList(c *gin.Context) { // change user's info
	// payload
	var payload entity.BatchRequestPayloadWithStringData
	if err := c.ShouldBindJSON(&payload); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// check permission
	tokenStr := c.GetHeader("Authorization")
	userSvc, _ := user.GetUserService()
	u, _ := userSvc.CheckToken(tokenStr)
	if u.GetRole() != "admin" && u.GetRole() != "root" {
		HandleErrorUnauthorized(c, errors.ErrorUserUnauthorized)
		return
	}
	// check whether its change root info
	// id, err := primitive.ObjectIDFromHex(c.Param("id"))
	// res, _ := c.Get(constants.ContextUser)
	// currentuser, _ := res.(interfaces.User)
	// u, ok = res.(interfaces.User)
	// if err != nil {
	// 	HandleErrorBadRequest(c, err)
	// 	return
	// }

	// doc to update
	var doc models.User
	if err := json.Unmarshal([]byte(payload.Data), &doc); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// query
	query := bson.M{
		"_id": bson.M{
			"$in": payload.Ids,
		},
	}
	println(query)
	// // update users
	if err := ctr.ctx.modelSvc.GetBaseService(interfaces.ModelIdUser).UpdateDoc(query, &doc, payload.Fields); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// update passwords
	if utils.Contains(payload.Fields, "password") {
		for _, id := range payload.Ids {
			if err := ctr.ctx.userSvc.ChangePassword(id, doc.Password); err != nil {
				trace.PrintError(err)
			}
		}
	}

	HandleSuccess(c)
}

func (ctr *userController) PutList(c *gin.Context) { //add new users
	// users
	var users []models.User
	if err := c.ShouldBindJSON(&users); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// // Check root user by zhizhong
	tokenStr := c.GetHeader("Authorization")
	userSvc, _ := user.GetUserService()
	u, _ := userSvc.CheckToken(tokenStr)

	if u.GetRole() != "root" {
		HandleErrorUnauthorized(c, errors.ErrorHttpUnauthorized)
		return
	}

	for _, u := range users {
		if err := ctr.ctx.userSvc.Create(&interfaces.UserCreateOptions{
			Username: u.Username,
			Password: u.Password,
			Email:    u.Email,
			Role:     u.Role,
		}); err != nil {
			trace.TraceError(err)
		}
	}

	HandleSuccess(c)
}

type userContext struct {
	modelSvc service.ModelService
	userSvc  interfaces.UserService
}

func (ctx *userContext) changePassword(c *gin.Context) {
	// Check that the request is valid
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// protect root role
	fmt.Println(id)
	aim_user, _ := ctx.modelSvc.GetUserById(id)
	// s, _ := aim_user.(interfaces.User)
	// aim_user, err = d.svc.modelSvc.GetUserById(id)
	fmt.Println(aim_user)
	// fmt.Println(s)
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// // Check root user by zhizhong
	tokenStr := c.GetHeader("Authorization")
	userSvc, _ := user.GetUserService()
	u, _ := userSvc.CheckToken(tokenStr)

	if aim_user.GetRole() == "root" && u.GetRole() != "root" {
		HandleErrorUnauthorized(c, errors.ErrorUserUnauthorized)
		return
	}
	if u.GetId() != id && u.GetRole() != "admin" && u.GetRole() != "root" {
		HandleErrorUnauthorized(c, errors.ErrorHttpUnauthorized)
		return
	}

	var payload map[string]string
	if err := c.ShouldBindJSON(&payload); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	password, ok := payload["password"]
	if !ok {
		HandleErrorBadRequest(c, errors.ErrorUserMissingRequiredFields)
		return
	}
	if len(password) < 5 {
		HandleErrorBadRequest(c, errors.ErrorUserInvalidPassword)
		return
	}
	if err := ctx.userSvc.ChangePassword(id, password); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccess(c)
}

func newUserContext() *userContext {
	// context
	ctx := &userContext{}

	// dependency injection
	c := dig.New()
	if err := c.Provide(service.NewService); err != nil {
		panic(err)
	}
	if err := c.Provide(user.ProvideGetUserService()); err != nil {
		panic(err)
	}
	if err := c.Invoke(func(
		modelSvc service.ModelService,
		userSvc interfaces.UserService,
	) {
		ctx.modelSvc = modelSvc
		ctx.userSvc = userSvc
	}); err != nil {
		panic(err)
	}

	return ctx
}

func newUserController() *userController {
	modelSvc, err := service.GetService()
	if err != nil {
		panic(err)
	}

	ctr := NewListPostActionControllerDelegate(ControllerIdUser, modelSvc.GetBaseService(interfaces.ModelIdUser), getUserActions())
	d := NewListPostActionControllerDelegate(ControllerIdUser, modelSvc.GetBaseService(interfaces.ModelIdUser), getUserActions())
	ctx := newUserContext()

	return &userController{
		ListActionControllerDelegate: *ctr,
		d:                            *d,
		ctx:                          ctx,
	}
}
