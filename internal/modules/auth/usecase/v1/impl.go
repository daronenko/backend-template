package usecase

import (
	"context"

	"emperror.dev/errors"
	"github.com/daronenko/backend-template/internal/app/config"
	"github.com/daronenko/backend-template/internal/models"
	"github.com/daronenko/backend-template/internal/modules/auth/repo/v1"
	"github.com/daronenko/backend-template/pkg/logger"
	"github.com/daronenko/backend-template/pkg/utils"
	"github.com/google/uuid"
)

const (
	prefix        = "api-auth"
	cacheDuration = 3600
)

// User usecase
type User struct {
	repo   repo.Repo
	cache  repo.Cache
	logger logger.Logger
	conf   *config.Config
}

func NewUser(repo repo.Repo, cache repo.Cache, logger logger.Logger, conf *config.Config) Usecase {
	return &User{repo, cache, logger, conf}
}

// Register new user
func (u *User) Register(ctx context.Context, user *models.User) (*models.UserWithToken, error) {
	existsUser, err := u.repo.GetByEmail(ctx, user)
	if existsUser != nil || err == nil {
		return nil, ErrUserExists
	}

	if err = user.PrepareCreate(); err != nil {
		return nil, errors.Wrap(err, "usecase.User.Register.user.PrepareCreate")
	}

	createdUser, err := u.repo.Create(ctx, user)
	if err != nil {
		return nil, MapRepoError(err)
	}

	token, err := utils.GenerateJWTToken(createdUser, u.conf)
	if err != nil {
		return nil, errors.Wrap(err, "usecase.User.Register.utils.GenerateJWTToken")
	}

	createdUser.SanitizePassword()
	return &models.UserWithToken{
		User:  createdUser,
		Token: token,
	}, nil
}

// Get access to user, session cookie send via cookie
func (u *User) Login(ctx context.Context, user *models.User) (*models.UserWithToken, error) {
	foundUser, err := u.repo.GetByEmail(ctx, user)
	if err != nil {
		return nil, MapRepoError(err)
	}

	if err = foundUser.ComparePasswords(user.Password); err != nil {
		return nil, ErrUnauthorized
	}

	token, err := utils.GenerateJWTToken(foundUser, u.conf)
	if err != nil {
		return nil, errors.Wrap(err, "usecase.User.Login.utils.GenerateJWTToken")
	}

	foundUser.SanitizePassword()
	return &models.UserWithToken{
		User:  foundUser,
		Token: token,
	}, nil
}

// Update user model
func (u *User) Update(ctx context.Context, user *models.User) (*models.User, error) {
	if err := user.PrepareUpdate(); err != nil {
		return nil, ErrMissingUserFields
	}

	updatedUser, err := u.repo.Update(ctx, user)
	if err != nil {
		return nil, MapRepoError(err)

	}

	if err = u.cache.Delete(ctx, user.ID); err != nil {
		return nil, MapRepoError(err)
	}

	updatedUser.SanitizePassword()
	return updatedUser, nil
}

// Delete user by user id
func (u *User) Delete(ctx context.Context, userID uuid.UUID) error {
	if err := u.repo.Delete(ctx, userID); err != nil {
		return MapRepoError(err)
	}

	if err := u.cache.Delete(ctx, userID); err != nil {
		return MapRepoError(err)
	}

	return nil
}

// Get user by user id
func (u *User) GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	cachedUser, err := u.cache.GetByID(ctx, userID)
	if err != nil {
		return nil, MapRepoError(err)
	}
	if cachedUser != nil {
		return cachedUser, nil
	}

	user, err := u.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, MapRepoError(err)
	}

	if err = u.cache.Set(ctx, user); err != nil {
		return nil, MapRepoError(err)
	}

	user.SanitizePassword()
	return user, nil
}

// Find users by username
func (u *User) FindByUsername(ctx context.Context, name string, query *utils.PaginationQuery) (*models.UsersList, error) {
	return u.repo.FindByName(ctx, name, query)
}

// Get users with pagination
func (u *User) GetUsers(ctx context.Context, pq *utils.PaginationQuery) (*models.UsersList, error) {
	return u.repo.GetUsers(ctx, pq)
}
