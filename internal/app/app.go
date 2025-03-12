package app

import (
	"log"
	"os"
	"time"

	"github.com/daronenko/backend-template/internal/app/appctx"
	"github.com/daronenko/backend-template/internal/app/config"
	"github.com/daronenko/backend-template/internal/infra"
	"github.com/daronenko/backend-template/internal/server"
	"github.com/daronenko/backend-template/internal/services/auth"
	"github.com/daronenko/backend-template/internal/services/meta"
	"github.com/daronenko/backend-template/internal/services/session"
	"github.com/daronenko/backend-template/pkg/logger"

	"go.uber.org/fx"
)

func Options(ctx appctx.Ctx, additionalOpts ...fx.Option) []fx.Option {
	conf, err := config.New(ctx)
	if err != nil {
		log.Printf("error: failed to parse config: %v\n", err)
		os.Exit(1)
	}

	logger.Setup(&conf.App.Logger)

	baseOpts := []fx.Option{
		fx.Supply(conf),

		fx.WithLogger(logger.Fx),

		infra.Module(),
		server.Module(),

		meta.Module(),
		auth.Module(),
		session.Module(),

		fx.StartTimeout(1 * time.Second),
	}

	return append(baseOpts, additionalOpts...)
}

func New(ctx appctx.Ctx, additionalOpts ...fx.Option) *fx.App {
	return fx.New(Options(ctx, additionalOpts...)...)
}
