package boltstore

import (
	"encoding/binary"
	"errors"
	"time"

	"github.com/RocketChat/statuspage/config"
	"github.com/RocketChat/statuspage/store"
	bolt "github.com/etcd-io/bbolt"
)

type boltStore struct {
	*bolt.DB
}

//New creates a new bolt store
func New() (store.Store, error) {
	if config.Config == nil {
		return nil, errors.New("configuration doesn't seem to exist")
	}

	db, err := bolt.Open(config.Config.DataPath, 0600, &bolt.Options{Timeout: 15 * time.Second})
	if err != nil {
		return nil, err
	}

	return &boltStore{db}, nil
}

//itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
