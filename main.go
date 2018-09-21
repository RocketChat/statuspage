package main

import (
	"flag"

	"github.com/RocketChat/statuspage/config"
	"github.com/RocketChat/statuspage/core"
	"github.com/RocketChat/statuspage/router"
)

func main() {
	configFile := flag.String("configFile", "config.yaml", "Config File full path. Defaults to current folder")

	flag.Parse()

	if err := config.Load(*configFile); err != nil {
		panic(err)
	}

	if err := core.TwistItUp(); err != nil {
		panic(err)
	}

	router.Start()
}
