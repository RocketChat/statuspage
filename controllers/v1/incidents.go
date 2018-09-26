package v1

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/RocketChat/statuspage/core"
	"github.com/RocketChat/statuspage/models"
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

	c.JSON(http.StatusOK, incidents)
}

//IncidentCreate creates the service, ensuring the database is correct
func IncidentCreate(c *gin.Context) {
	var incident models.Incident

	if err := c.BindJSON(&incident); err != nil {
		return
	}

	if err := core.CreateIncident(&incident); err != nil {
		internalErrorHandlerDetailed(c, err)
		return
	}

	c.JSON(http.StatusCreated, &incident)
}

//IncidentUpdate updates the service, ensuring the database is correct
func IncidentUpdate(c *gin.Context) {
	var incident models.Incident

	if err := c.BindJSON(&incident); err != nil {
		return
	}

	if err := core.UpdateIncident(&incident); err != nil {
		internalErrorHandlerDetailed(c, err)
		return
	}

	c.JSON(http.StatusOK, &incident)
}

//IncidentDelete removes the service, ensuring the database is correct
func IncidentDelete(c *gin.Context) {
	idParam := c.Param("id")

	if idParam == "" {
		internalErrorHandlerDetailed(c, errors.New("invalid incident id passed"))
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		internalErrorHandlerDetailed(c, err)
		return
	}

	if err := core.DeleteIncident(id); err != nil {
		internalErrorHandlerDetailed(c, err)
		return
	}

	c.Status(http.StatusOK)
}
