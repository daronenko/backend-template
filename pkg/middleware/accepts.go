package middleware

import (
	"strings"

	"github.com/daronenko/backend-template/internal/pkg/errs"
	"github.com/gofiber/fiber/v2"
)

func Accepts(mimes ...string) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		if ctx.Accepts(mimes...) != "" {
			return nil
		}

		return errs.ErrInvalidReq.Msg("invalid or missing 'accept' header, accepts: %s", strings.Join(mimes, ", "))
	}
}

var AcceptsJSON = Accepts("application/json")
