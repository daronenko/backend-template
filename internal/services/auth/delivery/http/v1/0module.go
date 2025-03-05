package httpdelivery

import (
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module("auth.delivery.v1", fx.Invoke(
		InitAuth,
	))
}
