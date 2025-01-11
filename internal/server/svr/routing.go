package svr

import (
	"crypto/subtle"
	"strings"

	"github.com/daronenko/backend-template/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type V1 struct {
	fiber.Router
}

type Admin struct {
	fiber.Router
}

type Meta struct {
	fiber.Router
}

func CreateEndpointGroups(app *fiber.App, cfg *config.Config) (*V1, *Admin, *Meta) {
	v1 := app.Group("/api/v1")

	admin := app.Group("/api/admin", func(c *fiber.Ctx) error {
		if len(cfg.Service.AdminKey) < 64 {
			log.Error().Msg("admin key is not set or is too short (at least should be 64 chars long), and a request has reached")
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		key := strings.TrimSpace(strings.TrimPrefix(c.Get(fiber.HeaderAuthorization), "Bearer"))

		// use constant time comparison to prevent timing attacks
		if subtle.ConstantTimeCompare([]byte(key), []byte(cfg.Service.AdminKey)) != 1 {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		return c.Next()
	})

	meta := app.Group("/api/_")

	return &V1{v1}, &Admin{admin}, &Meta{meta}
}
