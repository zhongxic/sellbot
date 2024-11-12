package main

import (
	"flag"
	"log"

	"github.com/zhongxic/sellbot/config"
	"github.com/zhongxic/sellbot/pkg/logger"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "config/config.yaml", "the config file in yaml format")
	flag.Parse()
}

func main() {
	cfg, err := config.Parse(configFile)
	if err != nil {
		log.Fatal(err)
	}
	if err := logger.Init(cfg.Logging); err != nil {
		log.Fatal(err)
	}
}
