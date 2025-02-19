package httpdelivery

import (
	"github.com/daronenko/backend-template/internal/pkg/errs"
	"github.com/gofiber/fiber/v2"
)

const (
	CodeServiceUnhealthy = "SERVICE_UNHEALTHY"
)

var (
	ErrServiceUnhealthy = errs.New(fiber.StatusInternalServerError, CodeServiceUnhealthy, "service is unhealthy")
)
