package middleware

import (
	"net/http"

	"exusiai.dev/gommon/constant"
	"github.com/getsentry/sentry-go"
	"github.com/gofiber/contrib/fibersentry"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

func EnrichSentry() func(ctx *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if hub := fibersentry.GetHubFromContext(c); hub != nil {
			hub.Scope().SetTag("request_id", c.Locals(constant.ContextKeyRequestID).(string))
		}

		var r http.Request
		if err := fasthttpadaptor.ConvertRequest(c.Context(), &r, true); err != nil {
			return err
		}

		spanIgnored := c.Get(constant.SlimHeaderKey) != ""

		if spanIgnored {
			return c.Next()
		}

		span := sentry.StartSpan(c.Context(), "backend", sentry.ContinueFromRequest(&r))
		span.SetTag("url", c.OriginalURL())
		span.Name = c.Method() + " " + c.Path()
		defer span.Finish()

		return c.Next()
	}
}
