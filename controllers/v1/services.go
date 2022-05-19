package v1

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/RocketChat/statuscentral/core"
	"github.com/RocketChat/statuscentral/models"
	"github.com/gin-gonic/gin"
)

// ServicesCreate creates a service
// @Summary Creates a service
// @ID services-create
// @Tags services
// @Accept json
// @Param service body models.Service true "Service object"
// @Produce json
// @Success 200 {object} models.Service
// @Router /v1/services [post]
func ServiceCreate(c *gin.Context) {
	var service models.Service

	if err := c.BindJSON(&service); err != nil {
		return
	}

	if err := core.CreateService(&service); err != nil {
		internalErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, service)
}

// ServicesGet gets all of the services
// @Summary Gets list of services
// @ID services-getall
// @Tags services
// @Produce json
// @Success 200 {object} []models.Service
// @Router /v1/services [get]
func ServicesGetAll(c *gin.Context) {
	services, err := core.GetServices()
	if err != nil {
		internalErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, services)
}

// ServicesGetOne gets one of the services
// @Summary Gets one of services
// @ID services-getone
// @Tags services
// @Produce json
// @Success 200 {object} models.Service
// @Router /v1/services/{id} [get]
func ServicesGetOne(c *gin.Context) {
	idParam := c.Param("id")

	if idParam == "" {
		badRequestHandlerDetailed(c, errors.New("invalid service id passed"))
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		badRequestHandlerDetailed(c, errors.New("invalid service id passed"))
		return
	}

	service, err := core.GetServiceByID(id)
	if err != nil {
		internalErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, service)
}

// ServicesGet gets all of the services
// @Summary Gets list of services
// @ID services-get
// @Tags services
// @Accept json
// @Param service body models.Service true "Service object"
// @Produce json
// @Success 200 {object} models.Service
// @Router /v1/services/{id} [post]
func ServiceUpdate(c *gin.Context) {
	idParam := c.Param("id")

	if idParam == "" {
		badRequestHandlerDetailed(c, errors.New("invalid service id passed"))
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		badRequestHandlerDetailed(c, errors.New("invalid service id passed"))
		return
	}

	var service models.Service

	if err := c.BindJSON(&service); err != nil {
		return
	}

	if id != service.ID {
		badRequestHandlerDetailed(c, errors.New("invalid service id passed"))
		return
	}

	if err := core.UpdateService(&service); err != nil {
		internalErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, service)
}
