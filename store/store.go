package store

import (
	"github.com/RocketChat/statuspage/models"
)

//Store is an interface that the storage implementers should implement
type Store interface {
	GetServices() ([]models.Service, error)
	GetIncidents(latest bool) ([]models.Incident, error)
	GetIncidentUpdates(id int64) ([]models.IncidentUpdate, error)
}
