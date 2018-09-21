package boltstore

import (
	"encoding/json"
	"time"

	"github.com/RocketChat/statuspage/models"
)

func (s *boltStore) GetServices() ([]*models.Service, error) {
	tx, err := s.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	cursor := tx.Bucket(incidentBucket).Cursor()

	services := make([]*models.Service, 0)
	for k, data := cursor.First(); k != nil; k, data = cursor.Next() {
		var i models.Service
		if err := json.Unmarshal(data, &i); err != nil {
			return nil, err
		}

		services = append(services, &i)
	}

	return services, nil
}

func (s *boltStore) CreateService(service *models.Service) error {
	tx, err := s.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	bucket := tx.Bucket(serviceBucket)

	seq, _ := bucket.NextSequence()
	service.ID = int(seq)
	service.UpdatedAt = time.Now()

	buf, err := json.Marshal(service)
	if err != nil {
		return err
	}

	if err := bucket.Put(itob(service.ID), buf); err != nil {
		return err
	}

	return tx.Commit()
}
