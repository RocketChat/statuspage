package core

import (
	"log"

	"github.com/RocketChat/statuscentral/store"
	"github.com/RocketChat/statuscentral/store/boltstore"
)

var _dataStore store.Store

// TwistItUp takes everything and starts the core up
func TwistItUp() error {
	store, err := boltstore.New()
	if err != nil {
		log.Fatalln(err)
	}

	_dataStore = store

	// Now that we have a store, let's ensure the services and regions from the config exist
	if err := createServicesFromConfig(); err != nil {
		log.Fatalln(err)
	}

	// Regions always need to be created AFTER services due to regions being tethered to services
	if err := createRegionsFromConfig(); err != nil {
		log.Fatalln(err)
	}

	return nil
}

// LivenessCheck checks the database to see if it responds to a ping
func LivenessCheck() error {
	return _dataStore.CheckDb()
}
