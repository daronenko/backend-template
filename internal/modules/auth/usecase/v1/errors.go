package usecase

import (
	"errors"

	"github.com/daronenko/backend-template/internal/modules/auth/repo/v1"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserExists        = errors.New("user with this username or email already exists")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrMissingUserFields = errors.New("some user model fields are missing")
)

func MapRepoError(err error) error {
	switch {
	case errors.Is(err, repo.ErrUserNotFound):
		return ErrUserNotFound
	case errors.Is(err, repo.ErrUserExists):
		return ErrUserExists
	default:
		return err
	}
}
