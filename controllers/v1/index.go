package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//IndexHandler is the html controller for sending the html dashboard
func IndexHandler(c *gin.Context) {
	// services := c.Keys["services"].(src.Services)
	// incidents := c.Keys["incidents"].(src.Incidents)

	// res, err := services.GetServices()
	// if err != nil {
	// 	panic(err)
	// }

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
		"owner":           "Rocket.Chat",
		"backgroundColor": "#343434",
		// "logo":               logo,
		// "services":           src.AggregateServices(res),
		// "mostCriticalStatus": src.MostCriticalStatus(res),
		// "incidents":          src.AggregateIncidents(inc),
	})
}
