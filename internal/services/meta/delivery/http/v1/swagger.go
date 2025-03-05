package httpdelivery

import (
	"github.com/daronenko/backend-template/docs"
	"github.com/daronenko/backend-template/internal/pkg/bininfo"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func InitSwagger(app *fiber.App) {
	docs.SwaggerInfo.Version = bininfo.Version
	app.Get("/swagger/*", swagger.HandlerDefault)
}
