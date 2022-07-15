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

// ScheduledMaintenanceGetAll gets all of the scheduled maintenance, latest depends on the "?all=true" query
// @Summary Gets list of incidents
// @ID scheduled-maintenance-getall
// @Tags scheduled-maintenance
// @Produce json
// @Success 200 {object} []models.ScheduledMaintenance
// @Router /v1/scheduled-maintenance [get]
func ScheduledMaintenanceGetAll(c *gin.Context) {
	allParam := c.Query("all")

	latest := true
	if allParam != "" && allParam == "true" {
		latest = false
	}

	scheduledMaintenance, err := core.GetScheduledMaintenance(latest)
	if err != nil {
		internalErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, scheduledMaintenance)
}

// ScheduledMaintenanceGetOne gets one scheduledMaintenance by the provided id
// @Summary Gets one scheduled maintenance
// @ID scheduled-maintenance-getOne
// @Tags scheduled-maintenance
// @Produce json
// @Success 200 {object} models.ScheduledMaintenance
// @Router /v1/scheduled-maintenance/{id} [get]
func ScheduledMaintenanceGetOne(c *gin.Context) {
	idParam := c.Param("id")

	if idParam == "" {
		badRequestHandlerDetailed(c, errors.New("invalid scheduled maintenance id passed"))
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		badRequestHandlerDetailed(c, err)
		return
	}

	incident, err := core.GetScheduledMaintenanceByID(id)
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

// ScheduledMaintenanceCreate creates a scheduled maintenance, ensuring the database is correct
// @Summary Creates a new scheduled maintenance
// @ID scheduled-maintenance-create
// @Tags scheduled-maintenance
// @Accept json
// @Param region body models.ScheduledMaintenance true "Scheduled Maintenance object"
// @Produce json
// @Success 200 {object} models.ScheduledMaintenance
// @Router /v1/scheduled-maintenance [post]
func ScheduledMaintenanceCreate(c *gin.Context) {
	var scheduledMaintenance models.ScheduledMaintenance

	if err := c.BindJSON(&scheduledMaintenance); err != nil {
		return
	}

	if scheduledMaintenance.Title == "" {
		badRequestHandlerDetailed(c, errors.New("title must be provided"))
		return
	}

	if scheduledMaintenance.PlannedStart.IsZero() {
		badRequestHandlerDetailed(c, errors.New("schedule maintenance must have a start date"))
		return
	}

	if scheduledMaintenance.PlannedEnd.IsZero() {
		badRequestHandlerDetailed(c, errors.New("schedule maintenance must have predicted end date"))
		return
	}

	maint, err := core.CreateScheduledMaintenance(&scheduledMaintenance)
	if err != nil {
		internalErrorHandlerDetailed(c, err)
		return
	}

	c.JSON(http.StatusCreated, &maint)
}

// ScheduledMaintenanceDelete removes the scheduled maintenance, ensuring the database is correct
// @Summary Deletes scheduled maintenance
// @ID scheduled-maintenance-delete
// @Tags scheduled-maintenance
// @Produce json
// @Success 200 {object} []models.ScheduleMaintenance
// @Router /v1/scheduled-maintenance/{id} [delete]
func ScheduledMaintenanceDelete(c *gin.Context) {
	idParam := c.Param("id")

	if idParam == "" {
		badRequestHandlerDetailed(c, errors.New("invalid scheduled maintenance id passed"))
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		badRequestHandlerDetailed(c, err)
		return
	}

	if err := core.DeleteScheduledMaintenance(id); err != nil {
		internalErrorHandlerDetailed(c, err)
		return
	}

	c.Status(http.StatusOK)
}

// ScheduledMaintenanceUpdateCreate creates an update for a scheduled maintenance
// @Summary Creates a scheduled maintenance update
// @ID scheduled-maintenance-create-update
// @Tags scheduled-maintenance
// @Accept json
// @Param region body models.IncidentUpdate true "Incident update object"
// @Param id path integer true "Incident id"
// @Produce json
// @Success 200 {object} models.IncidentUpdate
// @Router /v1/scheduled-maintenance/{id}/updates [post]
func ScheduledMaintenanceUpdateCreate(c *gin.Context) {
	idParam := c.Param("id")

	if idParam == "" {
		badRequestHandlerDetailed(c, errors.New("invalid scheduled maintenance id passed"))
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		badRequestHandlerDetailed(c, err)
		return
	}

	var update models.StatusUpdate
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

	maint, err := core.CreateScheduledMaintenanceUpdate(id, &update)
	if err != nil {
		internalErrorHandlerDetailed(c, err)
		return
	}

	c.JSON(http.StatusCreated, maint)
}

// ScheduledMaintenanceUpdatesGetAll gets updates for a scheduled maintenance
// @Summary Gets scheduled maintenance updates
// @ID scheduled-maintenance-update-getall
// @Tags scheduled-maintenance-update
// @Produce json
// @Success 200 {object} []models.IncidentUpdate
// @Router /v1/scheduled-maintenance/{id}/updates [get]
func ScheduledMaintenanceUpdatesGetAll(c *gin.Context) {
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

	updates, err := core.GetScheduledMaintenanceUpdates(id)
	if err != nil {
		internalErrorHandlerDetailed(c, err)
		return
	}

	c.JSON(http.StatusOK, updates)
}

// ScheduledMaintenanceUpdateGetOne gets an update for a scheduledMaintenance
// @Summary Gets one scheduled Maintenance update
// @ID scheduled-maintenance-update-getone
// @Tags scheduled-maintenance-update
// @Produce json
// @Success 200 {object} models.IncidentUpdate
// @Router /v1/scheduled-maintenance/{id}/updates/{updateId} [get]
func ScheduledMaintenanceUpdateGetOne(c *gin.Context) {
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

	update, err := core.GetScheduledMaintenanceUpdate(id, updateID)
	if err != nil {
		internalErrorHandlerDetailed(c, err)
		return
	}

	c.JSON(http.StatusOK, update)
}

// ScheduledMaintenanceUpdateDelete deletes an update for a scheduled maintenance
// @Summary Deletes one scheduled maintenance update
// @ID scheduled-maintenance-update-delete
// @Tags scheduled-maintenance-update
// @Produce json
// @Success 200 {object} models.IncidentUpdate
// @Router /v1/scheduled-maintenance/{id}/updates/{updateId} [delete]
func ScheduledMaintenanceUpdateDelete(c *gin.Context) {
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

	if err := core.DeleteScheduledMaintenanceUpdate(id, updateId); err != nil {
		internalErrorHandlerDetailed(c, err)
		return
	}

	c.Status(http.StatusOK)
}
