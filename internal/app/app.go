package app

import (
	"log"
	"time"

	"github.com/daronenko/backend-template/internal/config"
	"github.com/daronenko/backend-template/internal/delivery"
	"github.com/daronenko/backend-template/internal/server"
	"github.com/daronenko/backend-template/pkg/logger/fxlogger"
	"github.com/daronenko/backend-template/pkg/logger/zl"

	"go.uber.org/fx"
)

func Options(additionalOpts ...fx.Option) []fx.Option {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	// logger and configuration are the only two things that are not in the fx graph
	// because some other packages need them to be initialized before fx starts
	logger := zl.New(cfg.Service.Logger)

	baseOpts := []fx.Option{
		// fx meta
		fx.WithLogger(fxlogger.Fx),

		// Misc
		fx.Supply(cfg),
		fx.Supply(logger),

		// Servers
		server.Module(),

		// Delivery
		delivery.Module(delivery.OptIncludeSwagger),

		// fx Extra Options
		fx.StartTimeout(1 * time.Second),
		// StopTimeout is not typically needed, since we're using fiber's Shutdown(),
		// in which fiber has its own IdleTimeout for controlling the shutdown timeout.
		// It acts as a countermeasure in case the fiber app is not properly shutting down.
		fx.StopTimeout(5 * time.Minute),
	}

	return append(baseOpts, additionalOpts...)
}

func New(additionalOpts ...fx.Option) *fx.App {
	return fx.New(Options(additionalOpts...)...)
}
