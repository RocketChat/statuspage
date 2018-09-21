package core

import (
	"log"

	"github.com/RocketChat/statuspage/store"
	"github.com/RocketChat/statuspage/store/boltstore"
)

var _dataStore store.Store

//TwistItUp takes everything and starts the core up
func TwistItUp() error {
	log.Println("Creating the bolt data store")
	store, err := boltstore.New()
	if err != nil {
		log.Fatalln(err)
	}

	_dataStore = store

	return nil
}
