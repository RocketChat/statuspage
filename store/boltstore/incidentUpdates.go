package boltstore

import (
	"encoding/json"

	"github.com/RocketChat/statuspage/models"
)

func (s *boltStore) GetIncidentUpdates(id int64) ([]*models.IncidentUpdate, error) {
	tx, err := s.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	cursor := tx.Bucket(incidentBucket).Cursor()

	incidentUpdates := make([]*models.IncidentUpdate, 0)
	for k, data := cursor.First(); k != nil; k, data = cursor.Next() {
		var i models.IncidentUpdate
		if err := json.Unmarshal(data, &i); err != nil {
			return nil, err
		}

		incidentUpdates = append(incidentUpdates, &i)
	}

	return incidentUpdates, nil
}
