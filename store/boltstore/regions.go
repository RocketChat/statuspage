package boltstore

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/RocketChat/statuscentral/models"

	bolt "github.com/etcd-io/bbolt"
)

func (s *boltStore) GetRegions() ([]*models.Region, error) {
	tx, err := s.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	cursor := tx.Bucket(regionBucket).Cursor()

	regions := make([]*models.Region, 0)
	for k, data := cursor.First(); k != nil; k, data = cursor.Next() {
		var i models.Region
		if err := json.Unmarshal(data, &i); err != nil {
			return nil, err
		}

		regions = append(regions, &i)
	}

	return regions, nil
}

func (s *boltStore) GetRegionByName(name string) (*models.Region, error) {
	tx, err := s.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	cursor := tx.Bucket(regionBucket).Cursor()

	for k, data := cursor.First(); k != nil; k, data = cursor.Next() {
		var i models.Region
		if err := json.Unmarshal(data, &i); err != nil {
			return nil, err
		}

		if i.Name == name {
			return &i, nil
		}
	}

	return nil, nil
}

func (s *boltStore) UpdateRegion(region *models.Region) error {
	if region.ID <= 0 {
		return errors.New("invalid region id")
	}

	tx, err := s.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	bucket := tx.Bucket(regionBucket)

	region.UpdatedAt = time.Now()

	buf, err := json.Marshal(region)
	if err != nil {
		return err
	}

	if err := bucket.Put(itob(region.ID), buf); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *boltStore) CreateRegion(region *models.Region) error {
	tx, err := s.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	bucket := tx.Bucket(regionBucket)

	seq, _ := bucket.NextSequence()
	region.ID = int(seq)
	region.UpdatedAt = time.Now()

	buf, err := json.Marshal(region)
	if err != nil {
		return err
	}

	if err := bucket.Put(itob(region.ID), buf); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *boltStore) DeleteRegion(id int) error {
	return s.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(regionBucket).Delete(itob(id))
	})
}
