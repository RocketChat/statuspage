package core

import (
	"errors"

	"github.com/RocketChat/statuscentral/config"
	"github.com/RocketChat/statuscentral/models"
)

// GetRegions gets all of the services from the storage layer
func GetRegions() ([]*models.Region, error) {
	return _dataStore.GetRegions()
}

// GetRegionByCodeAndServiceName gets the region by code and service name, returns nil if not found
func GetRegionByCodeAndServiceName(regionCode, serviceName string) (*models.Region, error) {
	return _dataStore.GetRegionByCodeAndServiceName(regionCode, serviceName)
}

// CreateRegion creates the service in the storage layer
func CreateRegion(region *models.Region) error {
	return _dataStore.CreateRegion(region)
}

// ValidateAndCreateRegion checks if a region has all necessary info and creates it
func ValidateAndCreateRegion(region models.Region) (models.Region, error) {
	if region.Name == "" {
		return region, errors.New("region needs to have name")
	}

	if region.RegionCode == "" {
		return region, errors.New("region needs to have region code")
	}

	existingRegion, err := GetRegionByCodeAndServiceName(region.RegionCode, region.ServiceName)
	if err != nil {
		return region, err
	}

	if existingRegion != nil {
		return region, errors.New("region already exists")
	}

	service, err := GetServiceByName(region.ServiceName)
	if err != nil {
		return region, errors.New("could not find service with given serviceName")
	}

	region.ServiceID = service.ID

	if err := CreateRegion(&region); err != nil {
		return region, err
	}

	return region, nil
}

// DeleteRegion deletes a region in the storage layer
func DeleteRegion(id int) error {
	return _dataStore.DeleteRegion(id)
}

func createRegionsFromConfig() error {
	for _, pendingRegion := range config.Config.Regions {
		region, err := GetRegionByCodeAndServiceName(pendingRegion.RegionCode, pendingRegion.ServiceName)
		if err != nil {
			return err
		}

		if region != nil {
			continue
		}

		service, err := GetServiceByName(pendingRegion.ServiceName)
		if err != nil {
			return err
		}

		toCreate := &models.Region{
			Name:        pendingRegion.Name,
			Description: pendingRegion.Description,
			RegionCode:  pendingRegion.RegionCode,
			ServiceID:   service.ID,
			ServiceName: service.Name,
			Status:      models.ServiceStatusNominal,
			Enabled:     true,
			Tags:        make([]string, 0),
		}

		if err := CreateRegion(toCreate); err != nil {
			return err
		}
	}

	return nil
}

func updateRegionToStatus(regionCode, serviceName string, status models.ServiceAndRegionStatus) error {
	region, err := GetRegionByCodeAndServiceName(regionCode, serviceName)

	if err != nil {
		return err
	}

	if region == nil {
		return errors.New("unknown region")
	}

	val, ok := models.ServiceStatuses[status.ToLower()]

	if !ok {
		return errors.New("invalid region status")
	}

	region.Status = val

	return _dataStore.UpdateRegion(region)
}
