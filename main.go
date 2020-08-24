package main

import (
	"flag"

	"github.com/RocketChat/statuscentral/config"
	"github.com/RocketChat/statuscentral/core"
	"github.com/RocketChat/statuscentral/router"
)

func main() {
	configFile := flag.String("configFile", "statuscentral.yaml", "Config File full path. Defaults to current folder")

	flag.Parse()

	if err := config.Load(*configFile); err != nil {
		panic(err)
	}

	if err := core.TwistItUp(); err != nil {
		panic(err)
	}

	router.Start(config.Config.HTTP.Port)
}
