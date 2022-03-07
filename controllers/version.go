package controllers

import (
	"net/http"

	"github.com/buliqioqiolibusdo/demp-core/config"
	"github.com/gin-gonic/gin"
)

func GetVersion(c *gin.Context) {
	HandleSuccessWithData(c, config.GetVersion())
}

func getVersionActions() []Action {
	return []Action{
		{
			Method:      http.MethodGet,
			Path:        "",
			HandlerFunc: GetVersion,
		},
	}
}

var VersionController ActionController
