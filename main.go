package main

import (
	"flag"

	"github.com/RocketChat/statuscentral/config"
	"github.com/RocketChat/statuscentral/core"
	"github.com/RocketChat/statuscentral/router"
)

// @title           Status Central
// @version         0.1
// @description    	Operational status for Rocket Chat SaaS & Cloud

// @contact.name   Cloud Team
// @contact.url    https://open.rocket.chat/group/cloud

// @host status.rocket.chat
// @BasePath  /api
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
