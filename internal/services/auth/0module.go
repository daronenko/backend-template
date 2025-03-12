package auth

import (
	httpdelivery "github.com/daronenko/backend-template/internal/services/auth/delivery/http/v1"
	"github.com/daronenko/backend-template/internal/services/auth/repo/v1"
	"github.com/daronenko/backend-template/internal/services/auth/usecase/v1"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module("auth",
		repo.Module(),
		usecase.Module(),
		httpdelivery.Module(),
	)
}
