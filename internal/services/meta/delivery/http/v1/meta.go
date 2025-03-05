package httpdelivery

import (
	"time"

	"github.com/daronenko/backend-template/internal/pkg/bininfo"
	"github.com/daronenko/backend-template/internal/pkg/constant"
	"github.com/daronenko/backend-template/internal/server/svr"
	"github.com/daronenko/backend-template/internal/services/meta/usecase/v1"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"go.uber.org/fx"
)

type Meta struct {
	fx.In

	HealthUsecase usecase.HealthUsecase
}

func InitMeta(d Meta, meta *svr.Meta) {
	meta.Get("/bininfo", d.BinInfo)
	meta.Get("/ping", d.Ping)
	meta.Get("/health", cache.New(cache.Config{
		// cache it for a second to mitigate potential DDoS
		Expiration:  time.Second,
		CacheHeader: constant.CacheHeader,
	}), d.Health)
}

func (d *Meta) BinInfo(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"version":      bininfo.Version,
		"git_revision": bininfo.GitRevision,
	})
}

func (d *Meta) Ping(c *fiber.Ctx) error {
	// only allow intranet access to prevent abuse
	return c.SendString("pong")
}

func (d *Meta) Health(c *fiber.Ctx) error {
	if err := d.HealthUsecase.Ping(c.UserContext()); err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"status": "ok",
	})
}
