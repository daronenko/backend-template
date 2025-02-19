package repo

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New("user with this username or email already exists")
)
