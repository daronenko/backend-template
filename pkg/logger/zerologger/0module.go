package zerologger

import (
	. "github.com/daronenko/backend-template/pkg/logger/config"
	// . "github.com/daronenko/backend-template/pkg/logger/contracts"
	. "github.com/daronenko/backend-template/pkg/logger"
	"go.uber.org/fx"
)

func Module(cfg *Config) fx.Option {
	return fx.Module("zerologger", fx.Provide(
		func() Logger {
			return New(cfg)
		},
	))
}
