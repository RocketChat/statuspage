package core

import (
	"log"

	"github.com/RocketChat/statuscentral/store"
	"github.com/RocketChat/statuscentral/store/boltstore"
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

	//Now that we have a store, let's ensure the services from the config exist
	if err := createServicesFromConfig(); err != nil {
		log.Fatalln(err)
	}

	return nil
}
