package auth

import (
	httpdelivery "github.com/daronenko/backend-template/internal/modules/auth/delivery/http/v1"
	"github.com/daronenko/backend-template/internal/modules/auth/repo/v1"
	"github.com/daronenko/backend-template/internal/modules/auth/usecase/v1"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module("auth",
		repo.Module(),
		usecase.Module(),
		httpdelivery.Module(),
	)
}
