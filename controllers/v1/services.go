package v1

import (
	"github.com/RocketChat/statuspage/core"
	"github.com/gin-gonic/gin"
)

//ServicesGet gets all of the services
func ServicesGet(c *gin.Context) {
	services, err := core.GetServices()
	if err != nil {
		internalErrorHandler(c, err)
		return
	}

	c.JSON(200, services)
}
