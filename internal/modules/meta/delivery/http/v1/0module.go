package httpdelivery

import (
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module("meta.delivery.http.v1", fx.Invoke(
		RegisterMeta,
		RegisterIndex,
		RegisterSwagger,
	))
}
