package boltstore

import (
	"github.com/RocketChat/statuspage/models"
)

func (s *boltStore) GetIncidents(latest bool) ([]models.Incident, error) {
	return make([]models.Incident, 0), nil
}
