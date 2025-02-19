package usecase

import "go.uber.org/fx"

func Module() fx.Option {
	return fx.Module("meta.usecase.v1", fx.Provide(
		NewHealth,
	))
}
