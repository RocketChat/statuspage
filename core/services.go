package core

import (
	"errors"

	"github.com/RocketChat/statuscentral/config"
	"github.com/RocketChat/statuscentral/models"
)

// GetServices gets all of the services from the storage layer
func GetServices() ([]*models.Service, error) {
	return _dataStore.GetServices()
}

// GetServiceByName gets the service by name, returns nil if not found
func GetServiceByName(name string) (*models.Service, error) {
	return _dataStore.GetServiceByName(name)
}

// CreateService creates the service in the storage layer
func CreateService(service *models.Service) error {
	return _dataStore.CreateService(service)
}

// MostCriticalServiceStatus returns the most critical service number of the services provided
func MostCriticalServiceStatus(services []*models.Service, regions []*models.Region) int {
	mostCritical := 0

	for _, service := range services {
		serviceStatus := models.ServiceStatusValues[service.Status.String()]
		if serviceStatus > mostCritical {
			mostCritical = serviceStatus
		}
	}

	for _, region := range regions {
		regionStatus := models.ServiceStatusValues[region.Status.String()]
		if regionStatus > mostCritical {
			mostCritical = regionStatus
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
			Status:      models.ServiceStatusNominal,
			Enabled:     true,
			Tags:        make([]string, 0),
		}

		if err := CreateService(toCreate); err != nil {
			return err
		}
	}

	return nil
}

func updateServiceToStatus(serviceName string, status models.ServiceAndRegionStatus) error {
	service, err := GetServiceByName(serviceName)

	if err != nil {
		return err
	}

	if service == nil {
		return errors.New("unknown service")
	}

	val, ok := models.ServiceStatuses[status.ToLower()]

	if !ok {
		return errors.New("invalid service status")
	}

	service.Status = val

	return _dataStore.UpdateService(service)
}
