package httpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/daronenko/backend-template/internal/app/config"
	"github.com/daronenko/backend-template/internal/pkg/bininfo"
	"github.com/daronenko/backend-template/internal/pkg/errs"
	"github.com/daronenko/backend-template/internal/pkg/metrics"
	"github.com/daronenko/backend-template/pkg/middleware"
	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

var prometheusRegisterOnce sync.Once

type DevOpsApp struct {
	*fiber.App
}

func New(cfg *config.Config) (*fiber.App, DevOpsApp) {
	return CreateServiceApp(cfg), DevOpsApp{
		App: CreateDevOpsApp(cfg),
	}
}

func CreateServiceApp(conf *config.Config) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:               "Backend Template v1",
		ServerHeader:          fmt.Sprintf("BackendTemplate/%s", bininfo.Version),
		DisableStartupMessage: !conf.DevMode,

		// NOTICE: This will also affect WebSocket. Be aware if this fiber instance service is re-used
		//         for long connection services.
		ReadTimeout:  time.Second * 20,
		WriteTimeout: time.Second * 20,

		// allow possibility for graceful shutdown, otherwise app#Shutdown() will block forever
		IdleTimeout:             conf.Server.ShutdownTimeout,
		ProxyHeader:             "X-Original-Forwarded-For",
		EnableTrustedProxyCheck: true,
		TrustedProxies:          conf.Server.TrustedProxies,

		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,

		ErrorHandler: ErrorHandler,
		Immutable:    true,
	})

	app.Use(favicon.New())

	if !conf.DevMode {
		app.Use(cors.New(cors.Config{
			AllowOrigins: "*",
			AllowOriginsFunc: func(origin string) bool {
				return true
			},
			AllowMethods:     "GET, POST, DELETE, OPTIONS",
			AllowHeaders:     "Content-Type, Authorization, X-Requested-With, X-BackendTemplate-Variant",
			ExposeHeaders:    "Content-Type, X-BackendTemplate-Set-UserID, X-BackendTemplate-Upgrade, X-BackendTemplate-Compatible, X-BackendTemplate-Request-ID",
			AllowCredentials: true,
		}))
	}

	// requestid is used by report service to identify requests and generate taskId there afterwards
	// the logger middleware now injects RequestID into the context
	middleware.Logger(app)

	// then we need an extra middleware to extract it and repopulate it into ctx.Locals
	app.Use(middleware.RequestID())

	app.Use(func(c *fiber.Ctx) error {
		// use custom error handler to return customized error responses
		err := c.Next()
		if e, ok := err.(*errs.Error); ok {
			return HandleCustomError(c, e)
		}
		return err
	})

	app.Use(helmet.New(helmet.Config{
		HSTSMaxAge:         31356000,
		HSTSPreloadEnabled: true,
		ReferrerPolicy:     "strict-origin-when-cross-origin",
		PermissionPolicy:   "interest-cohort=()",
	}))

	// app.Use(middleware.InjectI18n())

	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e any) {
			buf := make([]byte, 4096)
			buf = buf[:runtime.Stack(buf, false)]
			log.Error().Msgf("panic: %v\n%s\n", e, buf)
		},
	}))

	tracerProvider := tracesdk.NewTracerProvider(
		append(
			[]tracesdk.TracerProviderOption{
				tracesdk.WithResource(resource.NewWithAttributes(
					semconv.SchemaURL,
					semconv.ServiceNameKey.String("backend"),
					semconv.ServiceVersionKey.String(bininfo.Version),
					semconv.ServiceInstanceIDKey.String(lo.Must(os.Hostname())),
					semconv.DeploymentEnvironmentKey.String(lo.Ternary(conf.DevMode, "dev", "prod")),
				)),
				tracesdk.WithSampler(
					tracesdk.ParentBased(
						tracesdk.TraceIDRatioBased(
							conf.App.Tracing.SampleRate))),
			},
			tracingProviderOptions(conf)...,
		)...,
	)
	otel.SetTracerProvider(tracerProvider)

	app.Use(otelfiber.Middleware(otelfiber.WithServerName("backend")))

	prometheusRegisterOnce.Do(func() {
		p := fiberprometheus.NewWithDefaultRegistry(metrics.Namespace)
		p.RegisterAt(app, "/metrics")
		app.Use(p.Middleware)
	})

	if conf.DevMode {
		log.Info().
			Str("evt.name", "infra.dev_mode.enabled").
			Msg("running in dev mode")
	}

	return app
}

func CreateDevOpsApp(conf *config.Config) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:               "Backend Template v1 (DevOps)",
		ServerHeader:          fmt.Sprintf("BackendTemplateDevOps/%s", bininfo.Version),
		DisableStartupMessage: !conf.DevMode,

		// allow possibility for graceful shutdown, otherwise app#Shutdown() will block forever
		IdleTimeout:             conf.Server.ShutdownTimeout,
		ProxyHeader:             "X-Original-Forwarded-For",
		EnableTrustedProxyCheck: true,
		TrustedProxies:          conf.Server.TrustedProxies,

		ErrorHandler: ErrorHandler,
		Immutable:    true,
	})

	app.Use(pprof.New())

	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e any) {
			buf := make([]byte, 4096)
			buf = buf[:runtime.Stack(buf, false)]
			log.Error().Msgf("panic: %v\n%s\n", e, buf)
		},
	}))

	return app
}

func tracingProviderOptions(conf *config.Config) []tracesdk.TracerProviderOption {
	options := []tracesdk.TracerProviderOption{}
	if !conf.App.Tracing.Enabled {
		log.Info().
			Str("evt.name", "infra.tracing.disabled").
			Msg("tracing is disabled: no spans will be reported")
		return options
	}

	optionsstr := make([]string, 0)

	if conf.App.Tracing.Exporters != nil {
		exporters := lo.Uniq(conf.App.Tracing.Exporters)
		for _, exporter := range exporters {
			switch exporter {
			case "otlpgrpc":
				exp := lo.Must(otlptracegrpc.New(context.Background()))
				options = append(options, tracesdk.WithBatcher(exp))
				optionsstr = append(optionsstr, "otlpgrpc")
			case "otlphttp":
				exp := lo.Must(otlptracehttp.New(
					context.Background(),
				))
				options = append(options, tracesdk.WithBatcher(exp))
				optionsstr = append(optionsstr, "otlphttp")
			case "stdout":
				exp := lo.Must(stdouttrace.New(stdouttrace.WithPrettyPrint()))
				options = append(options, tracesdk.WithSyncer(exp))
				optionsstr = append(optionsstr, "stdout")
			}
		}
	}

	if len(options) == 0 {
		log.Warn().
			Str("evt.name", "infra.tracing.exporters").
			Msg("tracing is enabled via configuration, but no tracing exporters are provided")
	} else {
		log.Info().
			Str("evt.name", "infra.tracing.exporters").
			Msgf("tracing enabled with exporters: %s", strings.Join(optionsstr, ", "))
	}

	return options
}
