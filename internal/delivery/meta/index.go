package meta

import "github.com/gofiber/fiber/v2"

func RegisterIndex(app *fiber.App) {
	app.Get("/api", Index)
}

func Index(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"@link":   "https://github.com/daronenko/backend-template",
		"message": "Welcome to Backend Template v1",
	})
}
