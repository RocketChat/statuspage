package v1

import (
	"net/http"

	"github.com/RocketChat/statuspage/config"
	"github.com/RocketChat/statuspage/core"
	"github.com/gin-gonic/gin"
)

//IndexHandler is the html controller for sending the html dashboard
func IndexHandler(c *gin.Context) {
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

	// logo := os.Getenv("SITE_LOGO")
	// if logo == "" {
	// 	logo = "static/img/logo.png"
	// }

	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"owner":           config.Config.Website.Title,
		"backgroundColor": config.Config.Website.HeaderBgColor,
		"logo":            "static/img/logo.svg",
		"services":        services,
		// "mostCriticalStatus": src.MostCriticalStatus(res),
		// "incidents":          src.AggregateIncidents(inc),
	})
}
