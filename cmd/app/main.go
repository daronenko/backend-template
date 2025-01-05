package main

import (
	"flag"
	"log"

	"github.com/daronenko/backend-template/internal/config"
	"github.com/daronenko/backend-template/internal/logger/zap"
)

var (
	Version  string
	Revision string
)

func main() {
	flag.Parse()

	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	logger := zap.New(cfg.Service.Logger)
	logger.Configure()
	logger.Info(cfg.Postgres)
}
