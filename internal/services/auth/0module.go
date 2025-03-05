package auth

import (
	httpdelivery "github.com/daronenko/backend-template/internal/services/auth/delivery/http/v1"
	"github.com/daronenko/backend-template/internal/services/auth/repo/v1"
	"github.com/daronenko/backend-template/internal/services/auth/usecase/v1"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

func Tracer() trace.Tracer {
	return otel.Tracer("auth")
}

func Module() fx.Option {
	return fx.Module("auth",
		fx.Provide(Tracer),
		repo.Module(),
		usecase.Module(),
		httpdelivery.Module(),
	)
}
