package store

import (
	"github.com/RocketChat/statuscentral/models"
)

//Store is an interface that the storage implementers should implement
type Store interface {
	GetServices() ([]*models.Service, error)
	GetIncidents(latest bool) ([]*models.Incident, error)

	GetServiceByName(name string) (*models.Service, error)
	GetIncidentByID(id int) (*models.Incident, error)

	CreateService(service *models.Service) error
	CreateIncident(incident *models.Incident) error

	CreateIncidentUpdate(incidentID int, update *models.IncidentUpdate) error
	GetIncidentUpdateByID(incidentID int, updateID int) (*models.IncidentUpdate, error)
	DeleteIncidentUpdateByID(incidentID int, updateID int) error

	UpdateService(service *models.Service) error
	UpdateIncident(incident *models.Incident) error

	DeleteIncident(id int) error

	CheckDb() error
}
