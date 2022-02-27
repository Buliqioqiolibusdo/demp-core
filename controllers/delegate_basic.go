package controllers

import (
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	delegate2 "github.com/crawlab-team/crawlab-core/models/delegate"
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
	res, _ := c.Get(constants.ContextUser)
	currentuser, _ := res.(interfaces.User)
	if currentuser.GetRole() != "root" {
		HandleErrorUnauthorized(c, errors.ErrorHttpUnauthorized)
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
	// Check whether the id of the operation object is its own id
	res, _ := c.Get(constants.ContextUser)
	currentuser, _ := res.(interfaces.User)
	if currentuser.GetId() != id && currentuser.GetRole() != "root" {
		HandleErrorUnauthorized(c, errors.ErrorHttpUnauthorized)
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
