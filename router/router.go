package router

import (
	"fmt"

	"github.com/RocketChat/statuscentral/config"
	v1c "github.com/RocketChat/statuscentral/controllers/v1"
	"github.com/RocketChat/statuscentral/router/middleware"
	"github.com/gin-gonic/gin"
)

//Start configures the routes and their handlers plus starts routing
func Start(port int) error {
	runMetricsRouter()

	router := gin.Default()

	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/*.tmpl")

	router.GET("/", v1c.IndexHandler)

	router.GET("/incidents/:id", v1c.IncidentDetailHandler)
	router.GET("/i/:id", v1c.IncidentShortRedirectHandler)

	v1 := router.Group("/api").Group("/v1")

	v1.GET("/services", v1c.ServicesGet)
	v1.GET("/incidents", v1c.IncidentsGetAll)
	v1.GET("/incidents/:id/updates", middleware.NotImplemented)

	v1.Use(middleware.IsAuthorized)
	{
		v1.GET("/config", config.Config.HttpHandler)

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

	return router.Run(fmt.Sprintf(":%d", port))
}

func runMetricsRouter() {
	healthMetricsRouter := gin.Default()
	healthMetricsRouter.GET("/health", v1c.LivenessCheckHandler)

	go healthMetricsRouter.Run(":8080")
}
