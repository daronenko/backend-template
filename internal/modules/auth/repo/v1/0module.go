package repo

import "go.uber.org/fx"

func Module() fx.Option {
	return fx.Module("auth.repo.v1", fx.Provide(
		NewUser,
		NewUserCache,
	))
}
