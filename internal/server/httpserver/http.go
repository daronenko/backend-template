package httpserver

import (
	"encoding/json"
	"fmt"
	"runtime"
	"time"

	"github.com/daronenko/backend-template/internal/app/config"
	"github.com/daronenko/backend-template/internal/pkg/bininfo"
	"github.com/daronenko/backend-template/internal/pkg/errs"
	"github.com/daronenko/backend-template/pkg/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog/log"
)

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

	// app.Use(fibersentry.New(fibersentry.Config{
	// 	Repanic: true,
	// 	Timeout: time.Second * 5,
	// }))

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowOriginsFunc: func(origin string) bool {
			return true
		},
		AllowMethods:     "GET, POST, DELETE, OPTIONS",
		AllowHeaders:     "Content-Type, Authorization, X-Requested-With, X-BackendTemplate-Variant, sentry-trace",
		ExposeHeaders:    "Content-Type, X-BackendTemplate-Set-UserID, X-BackendTemplate-Upgrade, X-BackendTemplate-Compatible, X-BackendTemplate-Request-ID",
		AllowCredentials: true,
	}))

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

	// tracerProvider := tracesdk.NewTracerProvider(
	// 	append(
	// 		[]tracesdk.TracerProviderOption{
	// 			tracesdk.WithResource(resource.NewWithAttributes(
	// 				semconv.SchemaURL,
	// 				semconv.ServiceNameKey.String("pgbackend"),
	// 				semconv.ServiceVersionKey.String(bininfo.Version),
	// 				semconv.ServiceInstanceIDKey.String(lo.Must(os.Hostname())),
	// 				semconv.DeploymentEnvironmentKey.String(lo.Ternary(conf.DevMode, "dev", "prod")),
	// 			)),
	// 			tracesdk.WithSampler(
	// 				tracesdk.ParentBased(
	// 					tracesdk.TraceIDRatioBased(
	// 						conf.TracingSampleRate))),
	// 		},
	// 		tracingProviderOptions(conf)...,
	// 	)...,
	// )
	// otel.SetTracerProvider(tracerProvider)

	// app.Use(otelfiber.Middleware(otelfiber.WithServerName("pgbackend")))

	// prometheusRegisterOnce.Do(func() {
	// 	fiberprometheus.New(observability.ServiceName).RegisterAt(app, "/metrics")
	// })

	if conf.DevMode {
		log.Info().
			Str("evt.name", "infra.dev_mode.enabled").
			Msg("running in dev mode")
	}

	// if !conf.DevMode {
	// 	app.Use(middleware.EnrichSentry())
	// }

	return app
}

func CreateDevOpsApp(conf *config.Config) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:               "Backend Template v3 (DevOps)",
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
