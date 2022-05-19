package v1

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/RocketChat/statuscentral/models"

	"github.com/RocketChat/statuscentral/config"
	"github.com/RocketChat/statuscentral/core"
	"github.com/gin-gonic/gin"
)

// IndexHandler is the html controller for sending the html dashboard
func IndexHandler(c *gin.Context) {
	services, err := core.GetServicesEnabled()
	if err != nil {
		log.Println("Error while getting the services:")
		log.Println(err)
		handleIndexPageLoadingFromConfig(c)
		return
	}

	incidents, err := core.GetIncidents(true)
	if err != nil {
		log.Println("Error while getting the incidents:")
		log.Println(err)
		handleIndexPageLoadingFromConfig(c)
		return
	}

	regions, err := core.GetRegions()
	if err != nil {
		log.Println("Error while getting the regions:")
		log.Println(err)
		handleIndexPageLoadingFromConfig(c)
		return
	}

	for _, service := range services {
		for _, region := range regions {
			if region.ServiceID == service.ID {
				service.Regions = append(service.Regions, *region)
			}
		}
	}

	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"owner":              config.Config.Website.Title,
		"backgroundColor":    config.Config.Website.HeaderBgColor,
		"cacheBreaker":       config.Config.Website.CacheBreaker,
		"logo":               "static/img/logo.svg",
		"services":           services,
		"mostCriticalStatus": core.MostCriticalServiceStatus(services, regions),
		"incidents":          core.AggregateIncidents(incidents),
	})
}

func handleIndexPageLoadingFromConfig(c *gin.Context) {
	services := make([]*models.Service, 0)
	for _, s := range config.Config.Services {
		service := &models.Service{
			Name:        s.Name,
			Description: s.Description,
			Status:      models.ServiceStatusUnknown,
		}

		services = append(services, service)
	}

	regions := make([]*models.Region, 0)
	for _, s := range config.Config.Regions {
		region := &models.Region{
			Name:        s.Name,
			Description: s.Description,
			Status:      models.ServiceStatusUnknown,
		}

		regions = append(regions, region)
	}

	for _, service := range services {
		for _, region := range regions {
			if region.ServiceID == service.ID {
				service.Regions = append(service.Regions, *region)
			}
		}
	}

	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"owner":              config.Config.Website.Title,
		"backgroundColor":    config.Config.Website.HeaderBgColor,
		"cacheBreaker":       config.Config.Website.CacheBreaker,
		"logo":               "static/img/logo.svg",
		"services":           services,
		"mostCriticalStatus": models.ServiceStatusValues["Unknown"],
		"incidents":          core.AggregateIncidents(make([]*models.Incident, 0)),
	})
}

func IncidentShortRedirectHandler(c *gin.Context) {
	if c.Param("id") == "" {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	c.Redirect(http.StatusPermanentRedirect, fmt.Sprintf("/incidents/%s", c.Param("id")))
}

// IncidentDetailHandler is the html controller for displaying the incident details
func IncidentDetailHandler(c *gin.Context) {
	if c.Param("id") == "" {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		internalErrorHandler(c, err)
		return
	}

	services, err := core.GetServices()
	if err != nil {
		log.Println("Error while getting the services:")
		log.Println(err)
		handleIndexPageLoadingFromConfig(c)
		return
	}

	regions, err := core.GetRegions()
	if err != nil {
		log.Println("Error while getting the regions:")
		log.Println(err)
		handleIndexPageLoadingFromConfig(c)
		return
	}

	for _, service := range services {
		for _, region := range regions {
			if region.ServiceID == service.ID {
				service.Regions = append(service.Regions, *region)
			}
		}
	}

	incident, err := core.GetIncidentByID(id)
	if err != nil {
		internalErrorHandler(c, err)
		return
	}

	if incident == nil {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	c.HTML(http.StatusOK, "detail.tmpl", gin.H{
		"owner":              config.Config.Website.Title,
		"backgroundColor":    config.Config.Website.HeaderBgColor,
		"cacheBreaker":       config.Config.Website.CacheBreaker,
		"logo":               "static/img/logo.svg",
		"mostCriticalStatus": core.MostCriticalServiceStatus(services, regions),
		"services":           services,
		"incident":           incident,
	})
}
