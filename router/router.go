package router

import (
	"net/http"

	"github.com/RocketChat/statuspage/config"
	"github.com/gin-gonic/gin"
)

func IndexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		/*"owner":              owner,
		"backgroundColor":    color,
		"logo":               logo,
		"services":           src.AggregateServices(res),
		"mostCriticalStatus": src.MostCriticalStatus(res),
		"incidents":          src.AggregateIncidents(inc),*/
	})
}

func NotImplemented(c *gin.Context) {
	c.AbortWithStatus(501)
}

// Just a temporary route
func ShowConfig(c *gin.Context) {
	c.JSON(200, config.Config)
}

func Start() {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/*")

	router.GET("/", IndexHandler)

	router.GET("/config", ShowConfig)

	api := router.Group("/api")
	{
		api.GET("/services", NotImplemented)
		api.POST("/services", NotImplemented)
		api.GET("/services/:id", NotImplemented)
		api.PATCH("/services/:id", NotImplemented)
		api.DELETE("/services/:id", NotImplemented)

		api.GET("/incidents", NotImplemented)
		api.POST("/incidents", NotImplemented)
		api.GET("/incidents/:id", NotImplemented)
		api.DELETE("/incidents/:id", NotImplemented)

		api.GET("/incidents/:id/updates", NotImplemented)
		api.POST("/incidents/:id/updates", NotImplemented)
		api.GET("/incidents/:id/updates/:updateId", NotImplemented)
		api.DELETE("/incidents/:id/updates/:updateId", NotImplemented)
	}

	router.Run(":5000")
}
