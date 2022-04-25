package v1

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/RocketChat/statuscentral/core"
	"github.com/RocketChat/statuscentral/models"
	"github.com/gin-gonic/gin"
)

// IncidentsGetAll gets all of the incidents, latest depends on the "?all=true" query
func IncidentsGetAll(c *gin.Context) {
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

// IncidentGetOne gets one incident by the provided id
func IncidentGetOne(c *gin.Context) {
	idParam := c.Param("id")

	if idParam == "" {
		badRequestHandlerDetailed(c, errors.New("invalid incident id passed"))
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		badRequestHandlerDetailed(c, err)
		return
	}

	incident, err := core.GetIncidentByID(id)
	if err != nil {
		internalErrorHandlerDetailed(c, err)
		return
	}

	if incident == nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, incident)
}

// IncidentCreate creates the incident, ensuring the database is correct
// @Summary Creates a new incident
// @ID incident-create
// @Tags incident
// @Accept json
// @Param region body models.Incident true "Incident object"
// @Produce json
// @Success 200 {object} models.Incident
// @Router /v1/incidents [post]
func IncidentCreate(c *gin.Context) {
	var incident models.Incident

	if err := c.BindJSON(&incident); err != nil {
		return
	}

	if incident.Title == "" {
		badRequestHandlerDetailed(c, errors.New("title must be provided"))
		return
	}

	if incident.Status == models.IncidentStatusScheduledMaintenance {
		if incident.Maintenance.Start.IsZero() {
			badRequestHandlerDetailed(c, errors.New("schedule maintenance incident must have a start date"))
			return
		}
		if incident.Maintenance.End.IsZero() {
			badRequestHandlerDetailed(c, errors.New("schedule maintenance incident must have predicted end date"))
			return
		}
	}

	inc, err := core.CreateIncident(&incident)
	if err != nil {
		internalErrorHandlerDetailed(c, err)
		return
	}

	c.JSON(http.StatusCreated, &inc)
}

// IncidentDelete removes the service, ensuring the database is correct
func IncidentDelete(c *gin.Context) {
	idParam := c.Param("id")

	if idParam == "" {
		badRequestHandlerDetailed(c, errors.New("invalid incident id passed"))
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		badRequestHandlerDetailed(c, err)
		return
	}

	if err := core.DeleteIncident(id); err != nil {
		internalErrorHandlerDetailed(c, err)
		return
	}

	c.Status(http.StatusOK)
}

// IncidentUpdateCreate creates an update for an incident
// @Summary Creates a new incident update
// @ID incident-create-update
// @Tags incident
// @Accept json
// @Param region body models.IncidentUpdate true "Incident update object"
// @Param id path integer true "Incident id"
// @Produce json
// @Success 200 {object} models.IncidentUpdate
// @Router /v1/incidents/{id}/updates [post]
func IncidentUpdateCreate(c *gin.Context) {
	idParam := c.Param("id")

	if idParam == "" {
		badRequestHandlerDetailed(c, errors.New("invalid incident id passed"))
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		badRequestHandlerDetailed(c, err)
		return
	}

	var update models.IncidentUpdate
	if err := c.BindJSON(&update); err != nil {
		return
	}

	if update.Message == "" {
		badRequestHandlerDetailed(c, errors.New("message is missing"))
		return
	}

	if update.Status == "" {
		badRequestHandlerDetailed(c, errors.New("status is missing"))
		return
	}

	status, ok := models.IncidentStatuses[strings.ToLower(update.Status.String())]
	if !ok {
		badRequestHandlerDetailed(c, errors.New("invalid status value"))
		return
	}

	update.Status = status

	incident, err := core.CreateIncidentUpdate(id, &update)
	if err != nil {
		internalErrorHandlerDetailed(c, err)
		return
	}

	c.JSON(http.StatusCreated, incident)
}

// IncidentUpdateGetOne gets an update for an incident
func IncidentUpdateGetOne(c *gin.Context) {
	idParam := c.Param("id")
	updateIdParam := c.Param("updateId")

	if idParam == "" {
		badRequestHandlerDetailed(c, errors.New("invalid incident id passed"))
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		badRequestHandlerDetailed(c, err)
		return
	}

	if updateIdParam == "" {
		badRequestHandlerDetailed(c, errors.New("invalid update id passed"))
		return
	}

	updateID, err := strconv.Atoi(updateIdParam)
	if err != nil {
		badRequestHandlerDetailed(c, err)
		return
	}

	update, err := core.GetIncidentUpdate(id, updateID)
	if err != nil {
		internalErrorHandlerDetailed(c, err)
		return
	}

	c.JSON(http.StatusOK, update)
}

// IncidentUpdateDelete deletes an update for an incident
func IncidentUpdateDelete(c *gin.Context) {
	idParam := c.Param("id")
	updateIdParam := c.Param("updateId")

	if idParam == "" {
		badRequestHandlerDetailed(c, errors.New("invalid incident id passed"))
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		badRequestHandlerDetailed(c, err)
		return
	}

	if updateIdParam == "" {
		badRequestHandlerDetailed(c, errors.New("invalid update id passed"))
		return
	}

	updateId, err := strconv.Atoi(updateIdParam)
	if err != nil {
		badRequestHandlerDetailed(c, err)
		return
	}

	if err := core.DeleteIncidentUpdate(id, updateId); err != nil {
		internalErrorHandlerDetailed(c, err)
		return
	}

	c.Status(http.StatusOK)
}
