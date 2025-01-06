package main

import (
	"flag"
	"log"
	"time"

	"github.com/pkg/errors"

	"github.com/daronenko/backend-template/internal/config"
	"github.com/daronenko/backend-template/internal/logger/zl"
)

func main() {
	flag.Parse()

	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	logger := zl.New(cfg.Service.Logger)
	err = errors.New("example error")
	logger.WarnErr("logging error", err)
	logger.ErrorErr("logging error", err)
	time.Sleep(180 * time.Second)
}
