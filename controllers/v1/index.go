package v1

import (
	"net/http"

	"github.com/RocketChat/statuspage/config"
	"github.com/RocketChat/statuspage/core"
	"github.com/gin-gonic/gin"
)

//IndexHandler is the html controller for sending the html dashboard
func IndexHandler(c *gin.Context) {
	// services := c.Keys["services"].(src.Services)
	// incidents := c.Keys["incidents"].(src.Incidents)

	services, err := core.GetServices()
	if err != nil {
		//TODO: fall back to the configuration
		internalErrorHandler(c, err)
		return
	}

	// inc, err := incidents.GetLatestIncidents()
	// if err != nil {
	// 	panic(err)
	// }

	// owner := os.Getenv("SITE_OWNER")
	// if owner == "" {
	// 	owner = "Abakus"
	// }

	// color := os.Getenv("SITE_COLOR")
	// if color == "" {
	// 	color = "#343434"
	// }

	// logo := os.Getenv("SITE_LOGO")
	// if logo == "" {
	// 	logo = "static/img/logo.png"
	// }

	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"owner":           config.Config.Website.Title,
		"backgroundColor": config.Config.Website.HeaderBgColor,
		"logo":            "static/img/logo.png",
		"services":        services,
		// "mostCriticalStatus": src.MostCriticalStatus(res),
		// "incidents":          src.AggregateIncidents(inc),
	})
}
