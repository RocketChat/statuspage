package boltstore

import (
	"encoding/json"
	"time"

	"github.com/RocketChat/statuspage/models"
)

func (s *boltStore) GetIncidents(latest bool) ([]*models.Incident, error) {
	tx, err := s.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	cursor := tx.Bucket(incidentBucket).Cursor()

	incidents := make([]*models.Incident, 0)
	for k, data := cursor.First(); k != nil; k, data = cursor.Next() {
		var i models.Incident
		if err := json.Unmarshal(data, &i); err != nil {
			return nil, err
		}

		incidents = append(incidents, &i)
	}

	return incidents, nil
}

func (s *boltStore) CreateIncident(incident *models.Incident) error {
	tx, err := s.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	bucket := tx.Bucket(incidentBucket)

	seq, _ := bucket.NextSequence()
	incident.ID = int(seq)
	incident.UpdatedAt = time.Now()

	buf, err := json.Marshal(incident)
	if err != nil {
		return err
	}

	if err := bucket.Put(itob(incident.ID), buf); err != nil {
		return err
	}

	return tx.Commit()
}
