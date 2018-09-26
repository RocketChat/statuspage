package boltstore

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/RocketChat/statuspage/config"
	"github.com/RocketChat/statuspage/models"
	bolt "github.com/etcd-io/bbolt"
)

func (s *boltStore) GetIncidents(latest bool) ([]*models.Incident, error) {
	tx, err := s.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	cursor := tx.Bucket(incidentBucket).Cursor()

	days := config.Config.Website.DaysToAggregate
	to := time.Now()
	from := to.Add(time.Duration(-days*24) * time.Hour).Truncate(24 * time.Hour)

	incidents := make([]*models.Incident, 0)
	for k, data := cursor.First(); k != nil; k, data = cursor.Next() {
		var i models.Incident
		if err := json.Unmarshal(data, &i); err != nil {
			return nil, err
		}

		if latest && i.Time.Before(to) && i.Time.After(from) {
			incidents = append(incidents, &i)
			continue
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

func (s *boltStore) UpdateIncident(incident *models.Incident) error {
	if incident.ID <= 0 {
		return errors.New("invalid incident id")
	}

	tx, err := s.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	bucket := tx.Bucket(incidentBucket)

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

func (s *boltStore) DeleteIncident(id int) error {
	return s.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(incidentBucket).Delete(itob(id))
	})
}
