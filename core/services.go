package core

import (
	"github.com/RocketChat/statuspage/config"
	"github.com/RocketChat/statuspage/models"
)

//GetServices gets all of the services from the storage layer
func GetServices() ([]*models.Service, error) {
	return _dataStore.GetServices()
}

//GetServiceByName gets the service by name, returns nil if not found
func GetServiceByName(name string) (*models.Service, error) {
	return _dataStore.GetServiceByName(name)
}

//CreateService creates the service in the storage layer
func CreateService(service *models.Service) error {
	return _dataStore.CreateService(service)
}

//MostCriticalServiceStatus returns the most critical service number of the services provided
func MostCriticalServiceStatus(services []*models.Service) int {
	mostCritical := 0

	for _, service := range services {
		serviceStatus := models.ServiceStatusValues[service.Status]
		if serviceStatus > mostCritical {
			mostCritical = serviceStatus
		}
	}

	return mostCritical
}

func createServicesFromConfig() error {
	for _, s := range config.Config.Services {
		service, err := GetServiceByName(s.Name)
		if err != nil {
			return err
		}

		if service != nil {
			continue
		}

		toCreate := &models.Service{
			Name:        s.Name,
			Description: s.Description,
			Status:      models.ServiceStatusOperational,
			Enabled:     true,
			Tags:        make([]string, 0),
		}

		if err := CreateService(toCreate); err != nil {
			return err
		}
	}

	return nil
}
