package repo

import "go.uber.org/fx"

func Module() fx.Option {
	return fx.Module("session.repo.v1", fx.Provide(
		NewSession,
	))
}
