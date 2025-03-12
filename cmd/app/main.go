package main

import (
	"context"
	"net"

	"github.com/daronenko/backend-template/internal/app"
	"github.com/daronenko/backend-template/internal/app/appctx"
	"github.com/daronenko/backend-template/internal/app/config"
	"github.com/daronenko/backend-template/internal/server/httpserver"
	"github.com/daronenko/backend-template/pkg/async"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
)

func main() {
	app.New(appctx.Declare(appctx.EnvServer), fx.Invoke(run)).Run()
}

func run(serviceApp *fiber.App, devOpsApp httpserver.DevOpsApp, conf *config.Config, lc fx.Lifecycle) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			serviceLn, err := net.Listen("tcp", conf.App.Address)
			if err != nil {
				return err
			}

			go func() {
				if err := serviceApp.Listener(serviceLn); err != nil {
					log.Error().Err(err).Msg("server terminated unexpectedly")
				}
			}()

			if conf.DevOps.Address == "" {
				log.Info().
					Str("evt.name", "infra.devops.disabled").
					Msg("devops server is disabled")
			} else {
				devOpsLn, err := net.Listen("tcp", conf.DevOps.Address)
				if err != nil {
					return err
				}

				go func() {
					if err := devOpsApp.Listener(devOpsLn); err != nil {
						log.Error().Err(err).Msg("server terminated unexpectedly")
					}
				}()
			}

			return nil
		},
		OnStop: func(ctx context.Context) error {
			if conf.DevMode {
				return nil
			}

			return async.WaitAll(
				async.Errable(serviceApp.Shutdown),
				async.Errable(devOpsApp.Shutdown),
			)
		},
	})
}
