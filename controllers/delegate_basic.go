package controllers

import (
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	delegate2 "github.com/crawlab-team/crawlab-core/models/delegate"
	"github.com/crawlab-team/crawlab-core/user"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
)

func NewBasicControllerDelegate(id ControllerId, svc interfaces.ModelBaseService) (d *BasicControllerDelegate) {
	return &BasicControllerDelegate{
		id:  id,
		svc: svc,
	}
}

type BasicControllerDelegate struct {
	id  ControllerId
	svc interfaces.ModelBaseService
}

func (d *BasicControllerDelegate) Get(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	doc, err := d.svc.GetById(id)
	if err == mongo2.ErrNoDocuments {
		HandleErrorNotFound(c, err)
		return
	}
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccessWithData(c, doc)
}

func (d *BasicControllerDelegate) Post(c *gin.Context) {
	// Check that the request is valid
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	// protect root role
	aim_user, _ := d.svc.GetById(id)
	s, _ := aim_user.(interfaces.User)
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	doc, err := NewJsonBinder(d.id).Bind(c)
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	if doc.GetId() != id {
		HandleErrorBadRequest(c, errors.ErrorHttpBadRequest)
		return
	}
	// Check whether the id of the operation object is its own id //testing by zhizhong
	tokenStr := c.GetHeader("Authorization")
	userSvc, _ := user.GetUserService()
	u, _ := userSvc.CheckToken(tokenStr)
	if s.GetRole() == "root" && u.GetRole() != "root" {
		HandleErrorUnauthorized(c, errors.ErrorUserUnauthorized)
		return
	}
	//fmt.Println(u.GetRole())
	if u.GetRole() != "root" && u.GetRole() != "admin" {
		HandleErrorUnauthorized(c, errors.ErrorUserUnauthorized)
		return
	}
	_, err = d.svc.GetById(id)
	if err != nil {
		HandleErrorNotFound(c, err)
		return
	}
	if err := delegate2.NewModelDelegate(doc, GetUserFromContext(c)).Save(); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccessWithData(c, doc)
}

func (d *BasicControllerDelegate) Put(c *gin.Context) {
	doc, err := NewJsonBinder(d.id).Bind(c)
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	if err := delegate2.NewModelDelegate(doc, GetUserFromContext(c)).Add(); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccessWithData(c, doc)
}

func (d *BasicControllerDelegate) Delete(c *gin.Context) {
	// Check that the request is valid
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// check permission
	tokenStr := c.GetHeader("Authorization")
	userSvc, _ := user.GetUserService()
	u, _ := userSvc.CheckToken(tokenStr)
	// fmt.Println(u.GetUsername())
	if u.GetRole() != "root" {
		// if u.GetUsername() != id && u.GetRole() != "root" && u.GetRole() != "admin" {
		HandleErrorUnauthorized(c, errors.ErrorUserUnauthorized)
		return
	}

	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	doc, err := d.svc.GetById(id)
	println(doc)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	if err := delegate2.NewModelDelegate(doc, GetUserFromContext(c)).Delete(); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccess(c)
}
