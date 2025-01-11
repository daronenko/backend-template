package meta

import (
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module("delivery.meta", fx.Invoke(
		RegisterMeta,
		RegisterIndex,
	))
}
