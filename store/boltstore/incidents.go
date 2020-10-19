package boltstore

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/RocketChat/statuscentral/config"
	"github.com/RocketChat/statuscentral/models"
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

func (s *boltStore) GetIncidentByID(id int) (*models.Incident, error) {
	tx, err := s.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	bytes := tx.Bucket(incidentBucket).Get(itob(id))
	if bytes == nil {
		return nil, nil
	}

	var i models.Incident
	if err := json.Unmarshal(bytes, &i); err != nil {
		return nil, err
	}

	return &i, nil
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

func (s *boltStore) CreateIncidentUpdate(incidentID int, update *models.IncidentUpdate) error {
	tx, err := s.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	bucket := tx.Bucket(incidentBucket)

	bytes := bucket.Get(itob(incidentID))
	if bytes == nil {
		return errors.New("no incident found by that id")
	}

	var i models.Incident
	if err := json.Unmarshal(bytes, &i); err != nil {
		return err
	}

	// If none index is 0 and then len should always put +1
	nextUpdateID := len(i.Updates)

	update.ID = nextUpdateID

	if update.Time.IsZero() {
		update.Time = time.Now()
	}

	i.Status = update.Status
	i.Updates = append(i.Updates, update)
	i.UpdatedAt = time.Now()

	buf, err := json.Marshal(i)
	if err != nil {
		return err
	}

	if err := bucket.Put(itob(i.ID), buf); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *boltStore) GetIncidentUpdateByID(incidentId int, updateId int) (*models.IncidentUpdate, error) {
	tx, err := s.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	bytes := tx.Bucket(incidentBucket).Get(itob(incidentId))
	if bytes == nil {
		return nil, nil
	}

	var incident models.Incident
	if err := json.Unmarshal(bytes, &incident); err != nil {
		return nil, err
	}

	for i, update := range incident.Updates {
		if update.ID == updateId {
			return incident.Updates[i], nil
		}
	}

	return nil, nil
}

func (s *boltStore) DeleteIncidentUpdateByID(incidentId int, updateId int) error {
	tx, err := s.Begin(false)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	bucket := tx.Bucket(incidentBucket)

	bytes := bucket.Get(itob(incidentId))
	if bytes == nil {
		return nil
	}

	var incident models.Incident
	if err := json.Unmarshal(bytes, &incident); err != nil {
		return err
	}

	updates := []*models.IncidentUpdate{}

	for _, update := range incident.Updates {
		if update.ID != updateId {
			updates = append(updates, update)
		}
	}

	incident.Updates = updates

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
