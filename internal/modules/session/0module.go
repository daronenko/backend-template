package session

import (
	"github.com/daronenko/backend-template/internal/modules/session/repo/v1"
	"github.com/daronenko/backend-template/internal/modules/session/usecase/v1"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module("session",
		repo.Module(),
		usecase.Module(),
	)
}
