package httpserver

import (
	"encoding/json"
	"fmt"
	"runtime"
	"time"

	"github.com/daronenko/backend-template/internal/config"
	"github.com/daronenko/backend-template/internal/pgerr"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
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

func CreateServiceApp(cfg *config.Config) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:               "Backend Template v1",
		ServerHeader:          fmt.Sprintf("Template/%s", config.Version),
		DisableStartupMessage: true,
		// NOTICE: This will also affect WebSocket. Be aware if this fiber instance service is re-used
		//         for long connection services.
		ReadTimeout:  time.Second * 20,
		WriteTimeout: time.Second * 20,
		// allow possibility for graceful shutdown, otherwise app#Shutdown() will block forever
		IdleTimeout:             time.Second * 60,
		ProxyHeader:             "X-Original-Forwarded-For",
		EnableTrustedProxyCheck: true,
		TrustedProxies:          []string{"::1", "127.0.0.1", "10.0.0.0/8"},

		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,

		ErrorHandler: ErrorHandler,
		Immutable:    true,
	})

	app.Use(favicon.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:8080,http://0.0.0.0:8080",
		AllowMethods:     "GET, POST, DELETE, OPTIONS",
		AllowHeaders:     "Content-Type, Authorization, X-Requested-With, X-Variant",
		ExposeHeaders:    "Content-Type, X-Request-ID",
		AllowCredentials: true,
	}))
	// requestid is used by report service to identify requests and generate taskId there afterwards
	// the logger middleware now injects RequestID into the context
	// middlewares.Logger(app)
	// then we need an extra middleware to extract it and repopulate it into ctx.Locals
	// app.Use(middlewares.RequestID())

	app.Use(func(c *fiber.Ctx) error {
		// Use custom error handler to return customized error responses
		err := c.Next()
		if e, ok := err.(*pgerr.PenguinError); ok {
			return HandleCustomError(c, e)
		}
		return err
	})

	// app.Use(middlewares.InjectI18n())
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

func CreateDevOpsApp(cfg *config.Config) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:               "Backend Template v1 (DevOps)",
		ServerHeader:          fmt.Sprintf("TemplateDevOps/%s", config.Version),
		DisableStartupMessage: false,
		// allow possibility for graceful shutdown, otherwise app#Shutdown() will block forever
		IdleTimeout:             time.Second * 60,
		ProxyHeader:             "X-Original-Forwarded-For",
		EnableTrustedProxyCheck: true,
		TrustedProxies:          []string{"::1", "127.0.0.1", "10.0.0.0/8"},

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
