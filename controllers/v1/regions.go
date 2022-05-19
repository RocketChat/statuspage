package v1

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/RocketChat/statuscentral/core"
	"github.com/RocketChat/statuscentral/models"
	"github.com/gin-gonic/gin"
)

// RegionCreate creates a new region
// @Summary Creates a new region
// @ID region-create
// @Tags region
// @Accept json
// @Param region body models.Region true "Region object"
// @Produce json
// @Success 200 {object} models.Region
// @Router /v1/regions [post]
func RegionCreate(c *gin.Context) {
	var region models.Region

	if err := c.BindJSON(&region); err != nil {
		return
	}

	region, err := core.ValidateAndCreateRegion(region)
	if err != nil {
		internalErrorHandlerDetailed(c, err)
		return
	}

	c.JSON(http.StatusCreated, region)
}

// RegionDelete deletes a given region
// @Summary Deletes a given region
// @ID region-delete
// @Tags region
// @Accept json
// @Param id path integer true "Region id"
// @Produce json
// @Success 204
// @Router /v1/regions/{id} [delete]
func RegionDelete(c *gin.Context) {
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

	if err := core.DeleteRegion(id); err != nil {
		internalErrorHandlerDetailed(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
