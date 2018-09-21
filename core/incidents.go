package core

import (
	"github.com/RocketChat/statuspage/models"
)

//GetIncidents retrieves the incidents from the storage layer
func GetIncidents(latest bool) ([]*models.Incident, error) {
	return _dataStore.GetIncidents(latest)
}
