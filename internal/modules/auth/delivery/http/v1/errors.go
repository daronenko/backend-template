package httpdelivery

import (
	"errors"

	"github.com/daronenko/backend-template/internal/modules/auth/usecase/v1"
	"github.com/daronenko/backend-template/internal/pkg/errs"
	"github.com/gofiber/fiber/v2"
)

const (
	CodeUserExists   = "USER_EXISTS"
	CodeUserNotFound = "USER_NOT_FOUND"
)

var (
	ErrUserExists           = errs.New(fiber.StatusBadRequest, CodeUserExists, "user with that name or email already exists")
	ErrUserNotFound         = errs.New(fiber.StatusNotFound, CodeUserNotFound, "user with that username and email does not exist")
	ErrMissingUsernameQuery = errs.New(fiber.StatusBadRequest, errs.CodeInvalidRequest, "username query param is required")
)

func MapUsecaseError(err error) *errs.Error {
	switch {
	case errors.Is(err, usecase.ErrUserExists):
		return ErrUserExists
	case errors.Is(err, usecase.ErrUserNotFound):
		return ErrUserNotFound
	case errors.Is(err, usecase.ErrMissingUserFields):
		return errs.ErrInvalidReq
	case errors.Is(err, usecase.ErrUnauthorized):
		return errs.ErrUnauthorized
	default:
		return errs.ErrInternalError
	}
}
