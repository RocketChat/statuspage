package main

import (
	"flag"

	"github.com/RocketChat/statuspage/config"
	"github.com/RocketChat/statuspage/router"
)

func main() {
	configFile := flag.String("configFile", "config.yaml", "Config File full path. Defaults to current folder")

	flag.Parse()

	err := config.Load(*configFile)
	if err != nil {
		panic(err)
	}

	router.Start()
}
