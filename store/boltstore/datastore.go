package boltstore

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/RocketChat/statuscentral/config"
	"github.com/RocketChat/statuscentral/store"
	bolt "github.com/etcd-io/bbolt"
)

type boltStore struct {
	*bolt.DB
}

var (
	incidentBucket = []byte("incidents")
	serviceBucket  = []byte("services")
)

//New creates a new bolt store
func New() (store.Store, error) {
	if config.Config == nil {
		return nil, errors.New("configuration doesn't seem to exist")
	}

	ex, err := os.Executable()
	if err != nil {
		return nil, err
	}

	db, err := bolt.Open(fmt.Sprintf("%s%s%s", filepath.Dir(ex), config.Config.DataPath, "statuscentral.bbolt"), 0600, &bolt.Options{Timeout: 15 * time.Second})
	if err != nil {
		return nil, err
	}

	tx, err := db.Begin(true)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if _, err := tx.CreateBucketIfNotExists(incidentBucket); err != nil {
		return nil, err
	}

	if _, err := tx.CreateBucketIfNotExists(serviceBucket); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &boltStore{db}, nil
}

func (s *boltStore) CheckDb() error {
	tx, err := s.Begin(false)
	if err != nil {
		return err
	}

	return tx.Rollback()
}

//itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
