package core

import (
	"github.com/RocketChat/statuspage/models"
)

//GetServices gets all of the services from the storage layer
func GetServices() ([]*models.Service, error) {
	return _dataStore.GetServices()
}
