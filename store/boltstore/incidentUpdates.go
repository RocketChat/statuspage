package boltstore

import (
	"github.com/RocketChat/statuspage/models"
)

func (s *boltStore) GetIncidentUpdates(id int64) ([]models.IncidentUpdate, error) {
	return make([]models.IncidentUpdate, 0), nil
}
