package main

import (
	"log"
	"time"

	"github.com/pkg/errors"

	"github.com/daronenko/backend-template/internal/config"
	"github.com/daronenko/backend-template/internal/logger/zl"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	logger := zl.New(cfg.Service.Logger)
	err = errors.New("example error")
	logger.WarnErr("logging error", err)
	logger.ErrorErr("logging error", err)
	logger.Debug(config.Version, config.Revision)
	time.Sleep(180 * time.Second)
}
