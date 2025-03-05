package httpdelivery

import (
	"fmt"

	"github.com/daronenko/backend-template/internal/pkg/bininfo"
	"github.com/gofiber/fiber/v2"
)

func InitIndex(app *fiber.App) {
	app.Get("/api", Index)
}

func Index(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"@link":   "https://github.com/daronenko/backend-template",
		"message": fmt.Sprintf("Welcome to Backend Template %s", bininfo.Version),
	})
}
