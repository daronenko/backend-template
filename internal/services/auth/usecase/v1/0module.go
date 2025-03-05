package usecase

import "go.uber.org/fx"

func Module() fx.Option {
	return fx.Module("auth.usecase.v1", fx.Provide(
		NewUser,
	))
}
