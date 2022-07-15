package client

import (
	"fmt"

	"github.com/RocketChat/statuscentral/models"
)

// ScheduledMaintenanceInterface ScheduledMaintenance interface
type ScheduledMaintenanceInterface interface {
	Create(scheduledMaintenance *models.ScheduledMaintenance) (returnedScheduledMaintenance *models.ScheduledMaintenance, err error)
	Get(id int) (scheduledMaintenance *models.ScheduledMaintenance, err error)
	GetMultiple(latestOnly bool) (result []*models.ScheduledMaintenance, err error)
	CreateStatusUpdate(maintenanceID int, statusUpdate *models.StatusUpdate) (returnedScheduledMaintenance *models.ScheduledMaintenance, err error)
	Delete(maintenanceID int) error
}

type scheduledMaintenance struct {
	client *Client
}

// Get Gets scheduled maintenance by id
func (i *scheduledMaintenance) Get(id int) (scheduledMaintenance *models.ScheduledMaintenance, err error) {

	req, err := i.client.buildRequest("GET", fmt.Sprintf("/api/v1/scheduled-maintenance/%d", id), nil)
	if err != nil {
		return nil, err
	}

	scheduledMaintenance = &models.ScheduledMaintenance{}

	resp, err := i.client.do(req, scheduledMaintenance)
	if err != nil {
		return nil, err
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return scheduledMaintenance, nil
}

func (i *scheduledMaintenance) GetMultiple(latestOnly bool) (result []*models.ScheduledMaintenance, err error) {
	req, err := i.client.buildRequest(
		"GET",
		fmt.Sprintf(
			"/api/v1/scheduled-maintenance?all=%v",
			!latestOnly,
		),
		nil,
	)

	if err != nil {
		return nil, err
	}

	result = []*models.ScheduledMaintenance{}

	resp, err := i.client.do(req, &result)
	if err != nil {
		return nil, err
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Create creates a scheduled maintenance
func (i *scheduledMaintenance) Create(scheduledMaintenance *models.ScheduledMaintenance) (returnedScheduledMaintenance *models.ScheduledMaintenance, err error) {
	req, err := i.client.buildRequest("POST", "/api/v1/scheduled-maintenance", scheduledMaintenance)
	if err != nil {
		return nil, err
	}

	returnedScheduledMaintenance = &models.ScheduledMaintenance{}

	resp, err := i.client.do(req, returnedScheduledMaintenance)
	if err != nil {
		return nil, err
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return returnedScheduledMaintenance, nil
}

// CreateStatusUpdate Creates a status update for a scheduled maintenance
func (i *scheduledMaintenance) CreateStatusUpdate(maintenanceID int, statusUpdate *models.StatusUpdate) (returnedScheduledMaintenance *models.ScheduledMaintenance, err error) {
	req, err := i.client.buildRequest("POST", fmt.Sprintf("/api/v1/scheduled-maintenance/%d/updates", maintenanceID), statusUpdate)
	if err != nil {
		return nil, err
	}

	returnedScheduledMaintenance = &models.ScheduledMaintenance{}

	resp, err := i.client.do(req, returnedScheduledMaintenance)
	if err != nil {
		return nil, err
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return returnedScheduledMaintenance, nil
}

// Delete deletes an incident
func (i *scheduledMaintenance) Delete(maintenanceID int) error {
	req, err := i.client.buildRequest(
		"DELETE",
		fmt.Sprintf("/api/v1/scheduled-maintenance/%d", maintenanceID),
		nil,
	)

	if err != nil {
		return err
	}

	resp, err := i.client.do(req, nil)
	if err != nil {
		return err
	}

	err = resp.Body.Close()

	return err
}
