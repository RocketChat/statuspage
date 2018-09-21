package boltstore

import "github.com/RocketChat/statuspage/models"

func (s *boltStore) GetServices() ([]models.Service, error) {
	return make([]models.Service, 0), nil
}
