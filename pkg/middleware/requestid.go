package middleware

import (
	"github.com/daronenko/backend-template/internal/pkg/constant"
	"github.com/daronenko/backend-template/pkg/flog"
	"github.com/gofiber/fiber/v2"
)

func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, ok := flog.IDFromFiberCtx(c)
		if ok {
			c.Locals(constant.ContextKeyRequestID, id.String())
		}
		return c.Next()
	}
}
