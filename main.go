package main

import (
	"flag"
	"fmt"
	"github.com/jetrmm/rmm-nats-api/api"

	"github.com/sirupsen/logrus"
)

var (
	version = "0.1.0"
	log     = logrus.New()
)

func main() {
	ver := flag.Bool("version", false, "Prints version")
	cfg := flag.String("config", "", "Path to the configuration file")
	logLevel := flag.String("log", "INFO", "The logging level")
	flag.Parse()

	if *ver {
		fmt.Println(version)
		return
	}

	setupLogging(logLevel)

	api.Svc(log, *cfg)
}

func setupLogging(level *string) {
	ll, err := logrus.ParseLevel(*level)
	if err != nil {
		ll = logrus.InfoLevel
	}
	log.SetLevel(ll)
}
