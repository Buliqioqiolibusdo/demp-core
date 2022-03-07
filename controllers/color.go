package controllers

import (
	"net/http"

	"github.com/buliqioqiolibusdo/demp-core/errors"
	"github.com/gin-gonic/gin"
)

func GetColorList(c *gin.Context) {
	panic(errors.ErrorControllerNotImplemented)
}

func getColorActions() []Action {
	return []Action{
		{Method: http.MethodGet, Path: "", HandlerFunc: GetColorList},
	}
}

var ColorController ActionController
