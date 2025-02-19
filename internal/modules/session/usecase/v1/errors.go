package usecase

import (
	"errors"

	"github.com/daronenko/backend-template/internal/modules/session/repo/v1"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New("user with this username or email already exists")
)

func MapRepoError(err error) error {
	switch {
	case errors.Is(err, repo.ErrUserNotFound):
		return ErrUserNotFound
	case errors.Is(err, repo.ErrUserExists):
		return ErrUserExists
	default:
		return nil
	}
}
