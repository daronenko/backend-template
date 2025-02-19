package app

import (
	"log"
	"os"
	"time"

	"github.com/daronenko/backend-template/internal/app/config"
	"github.com/daronenko/backend-template/internal/app/ctx"
	"github.com/daronenko/backend-template/internal/infra"
	"github.com/daronenko/backend-template/internal/modules/auth"
	"github.com/daronenko/backend-template/internal/modules/meta"
	"github.com/daronenko/backend-template/internal/modules/session"
	"github.com/daronenko/backend-template/internal/server"
	"github.com/daronenko/backend-template/pkg/logger/fxlogger"
	"github.com/daronenko/backend-template/pkg/logger/zerologger"

	"go.uber.org/fx"
)

func Options(ctx ctx.Ctx, additionalOpts ...fx.Option) []fx.Option {
	cfg, err := config.New(ctx)
	if err != nil {
		log.Printf("error: failed to parse config: %v\n", err)
		os.Exit(1)
	}

	baseOpts := []fx.Option{
		fx.Supply(cfg),

		fx.WithLogger(fxlogger.Fx),
		zerologger.Module(&cfg.Service.Logger),

		infra.Module(),
		server.Module(),

		meta.Module(),
		auth.Module(),
		session.Module(),

		fx.StartTimeout(1 * time.Second),
	}

	return append(baseOpts, additionalOpts...)
}

func New(ctx ctx.Ctx, additionalOpts ...fx.Option) *fx.App {
	return fx.New(Options(ctx, additionalOpts...)...)
}
