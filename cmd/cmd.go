package main

import (
	"flag"

	"github.com/supermatt/query-exporter-go/src/config"
	"github.com/supermatt/query-exporter-go/src/server"
)

func main() {
	configFile := flag.String("config", "", "config file")
	debug := flag.Bool("debug", false, "debug mode")
	flag.Parse()

	cfg, err := config.LoadConfig(*configFile)
	if err != nil {
		panic(err)
	}

	if *debug {
		cfg.Debug = true
	}

	s := server.NewServer(cfg)
	s.Start()
}
