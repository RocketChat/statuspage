package v1

import (
	"github.com/RocketChat/statuspage/core"
	"github.com/gin-gonic/gin"
)

//IncidentsGet gets the services, latest depends on the "?all=true" query
func IncidentsGet(c *gin.Context) {
	allParam := c.Query("all")

	latest := true
	if allParam != "" && allParam == "true" {
		latest = false
	}

	incidents, err := core.GetIncidents(latest)
	if err != nil {
		internalErrorHandler(c, err)
		return
	}

	c.JSON(200, incidents)
}
