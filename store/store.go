package store

import (
	"github.com/RocketChat/statuspage/models"
)

//Store is an interface that the storage implementers should implement
type Store interface {
	GetServices() ([]*models.Service, error)
	GetIncidents(latest bool) ([]*models.Incident, error)
	GetIncidentUpdates(id int64) ([]*models.IncidentUpdate, error)

	GetServiceByName(name string) (*models.Service, error)

	CreateService(service *models.Service) error
	CreateIncident(incident *models.Incident) error

	UpdateIncident(incident *models.Incident) error

	DeleteIncident(id int) error
}
