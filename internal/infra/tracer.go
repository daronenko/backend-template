package infra

import (
	"github.com/daronenko/backend-template/internal/app/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

func NewTracer(conf *config.Config) trace.Tracer {
	if !conf.App.Tracing.Enabled {
		return noop.NewTracerProvider().Tracer("")
	}

	return otel.Tracer("backend")
}
