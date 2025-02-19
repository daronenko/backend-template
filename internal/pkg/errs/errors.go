package errs

import (
	"github.com/gofiber/fiber/v2"
)

const (
	CodeNotFound       = "NOT_FOUND"
	CodeInvalidRequest = "INVALID_REQUEST"
	CodeInternalError  = "INTERNAL_ERROR"
	CodeUnauthorized   = "UNAUTHORIZED"
)

var (
	ErrNotFound               = New(fiber.StatusBadRequest, CodeNotFound, "resource not found with given parameters")
	ErrInvalidReq             = New(fiber.StatusBadRequest, CodeInvalidRequest, "invalid request: some or all request parameters are invalid")
	ErrInternalError          = New(fiber.StatusInternalServerError, CodeInternalError, "internal server error occurred")
	ErrInternalErrorImmutable = NewImmutable(fiber.StatusInternalServerError, CodeInternalError, "internal server error occurred")
	ErrUnauthorized           = New(fiber.StatusUnauthorized, CodeUnauthorized, "unauthorized")
)
