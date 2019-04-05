package boltstore

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/RocketChat/statuscentral/models"
)

func (s *boltStore) GetServices() ([]*models.Service, error) {
	tx, err := s.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	cursor := tx.Bucket(serviceBucket).Cursor()

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

func (s *boltStore) GetServiceByName(name string) (*models.Service, error) {
	tx, err := s.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	cursor := tx.Bucket(serviceBucket).Cursor()

	for k, data := cursor.First(); k != nil; k, data = cursor.Next() {
		var i models.Service
		if err := json.Unmarshal(data, &i); err != nil {
			return nil, err
		}

		if i.Name == name {
			return &i, nil
		}
	}

	return nil, nil
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

func (s *boltStore) UpdateService(service *models.Service) error {
	if service.ID <= 0 {
		return errors.New("invalid service id")
	}

	tx, err := s.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	bucket := tx.Bucket(serviceBucket)

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
