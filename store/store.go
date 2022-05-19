package store

import (
	"github.com/RocketChat/statuscentral/models"
)

//Store is an interface that the storage implementers should implement
type Store interface {
	// Services
	CreateService(service *models.Service) error
	UpdateService(service *models.Service) error
	GetServices() ([]*models.Service, error)
	GetServicesEnabled() ([]*models.Service, error)
	GetServiceByName(name string) (*models.Service, error)
	GetServiceByID(id int) (*models.Service, error)

	// Regions
	CreateRegion(region *models.Region) error
	UpdateRegion(region *models.Region) error
	GetRegions() ([]*models.Region, error)
	GetRegionByName(name string) (*models.Region, error)
	DeleteRegion(id int) error

	// Incidents
	CreateIncident(incident *models.Incident) error
	UpdateIncident(incident *models.Incident) error
	GetIncidents(latest bool) ([]*models.Incident, error)
	GetIncidentByID(id int) (*models.Incident, error)
	DeleteIncident(id int) error

	// Incident Updates
	CreateIncidentUpdate(incidentID int, update *models.IncidentUpdate) error
	GetIncidentUpdateByID(incidentID int, updateID int) (*models.IncidentUpdate, error)
	DeleteIncidentUpdateByID(incidentID int, updateID int) error

	CheckDb() error
}
