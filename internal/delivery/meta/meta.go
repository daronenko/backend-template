package meta

import (
	"github.com/daronenko/backend-template/internal/server/svr"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

type Meta struct {
	fx.In
}

func RegisterMeta(meta *svr.Meta, c Meta) {
	meta.Get("/health", c.Health)
}

func (c *Meta) Health(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"status": "ok",
	})
}
