package usecase

import "go.uber.org/fx"

func Module() fx.Option {
	return fx.Module("session.usecase.v1", fx.Provide(
		NewSession,
	))
}
