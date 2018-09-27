package router

import (
	"github.com/RocketChat/statuscentral/config"
	v1c "github.com/RocketChat/statuscentral/controllers/v1"
	"github.com/RocketChat/statuscentral/router/middleware"
	"github.com/gin-gonic/gin"
)

//ShowConfig is just a temporary route
func ShowConfig(c *gin.Context) {
	c.JSON(200, config.Config)
}

//Start configures the routes and their handlers plus starts routing
func Start() error {
	router := gin.Default()

	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/*")

	router.GET("/", v1c.IndexHandler)

	router.GET("/config", ShowConfig)

	v1 := router.Group("/api").Group("/v1")

	v1.GET("/services", v1c.ServicesGet)
	v1.GET("/incidents", v1c.IncidentsGetAll)
	v1.GET("/incidents/:id/updates", middleware.NotImplemented)

	v1.Use(middleware.IsAuthorized)
	{
		v1.POST("/services", middleware.NotImplemented)
		v1.GET("/services/:id", middleware.NotImplemented)
		v1.PATCH("/services/:id", middleware.NotImplemented)
		v1.DELETE("/services/:id", middleware.NotImplemented)

		v1.POST("/incidents", v1c.IncidentCreate)
		v1.GET("/incidents/:id", v1c.IncidentGetOne)
		v1.DELETE("/incidents/:id", v1c.IncidentDelete)

		v1.POST("/incidents/:id/updates", v1c.IncidentUpdateCreate)
		v1.GET("/incidents/:id/updates/:updateId", middleware.NotImplemented)
		v1.DELETE("/incidents/:id/updates/:updateId", middleware.NotImplemented)
	}

	return router.Run(":5000")
}
