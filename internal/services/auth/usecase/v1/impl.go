package usecase

import (
	"context"

	"emperror.dev/errors"
	"github.com/daronenko/backend-template/internal/app/config"
	"github.com/daronenko/backend-template/internal/model/v1"
	"github.com/daronenko/backend-template/internal/services/auth/repo/v1"
	"github.com/daronenko/backend-template/internal/util"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

// User usecase
type User struct {
	fx.In

	Conf   *config.Config
	Repo   repo.Repo
	Cache  repo.Cache
	Tracer trace.Tracer
}

// Needs to provide to fx
func NewUser(u User) Usecase {
	return &u
}

// Register new user
func (u *User) Register(ctx context.Context, user *model.User) (*model.UserWithToken, error) {
	existsUser, err := u.Repo.GetByEmail(ctx, user)
	if existsUser != nil || err == nil {
		return nil, ErrUserExists
	}

	if err = user.PrepareCreate(); err != nil {
		return nil, errors.Wrap(err, "usecase.User.Register.user.PrepareCreate")
	}

	createdUser, err := u.Repo.Create(ctx, user)
	if err != nil {
		return nil, MapRepoError(err)
	}

	token, err := util.GenerateJWTToken(createdUser, u.Conf)
	if err != nil {
		return nil, errors.Wrap(err, "usecase.User.Register.util.GenerateJWTToken")
	}

	createdUser.SanitizePassword()
	return &model.UserWithToken{
		User:  createdUser,
		Token: token,
	}, nil
}

// Get access to user, session cookie send via cookie
func (u *User) Login(ctx context.Context, user *model.User) (*model.UserWithToken, error) {
	ctx, span := u.Tracer.Start(ctx, "usercase.User.Login")
	defer span.End()

	foundUser, err := u.Repo.GetByEmail(ctx, user)
	if err != nil {
		return nil, MapRepoError(err)
	}

	if err = foundUser.ComparePasswords(user.Password); err != nil {
		return nil, ErrUnauthorized
	}

	token, err := util.GenerateJWTToken(foundUser, u.Conf)
	if err != nil {
		return nil, errors.Wrap(err, "usecase.User.Login.util.GenerateJWTToken")
	}

	foundUser.SanitizePassword()
	return &model.UserWithToken{
		User:  foundUser,
		Token: token,
	}, nil
}

// Update user model
func (u *User) Update(ctx context.Context, user *model.User) (*model.User, error) {
	if err := user.PrepareUpdate(); err != nil {
		return nil, ErrMissingUserFields
	}

	updatedUser, err := u.Repo.Update(ctx, user)
	if err != nil {
		return nil, MapRepoError(err)
	}

	if err = u.Cache.Delete(ctx, user.ID); err != nil {
		return nil, MapRepoError(err)
	}

	updatedUser.SanitizePassword()
	return updatedUser, nil
}

// Delete user by user id
func (u *User) Delete(ctx context.Context, userID uuid.UUID) error {
	if err := u.Repo.Delete(ctx, userID); err != nil {
		return MapRepoError(err)
	}

	if err := u.Cache.Delete(ctx, userID); err != nil {
		return MapRepoError(err)
	}

	return nil
}

// Get user by user id
func (u *User) GetByID(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	cachedUser, err := u.Cache.GetByID(ctx, userID)
	if err != nil {
		return nil, MapRepoError(err)
	}
	if cachedUser != nil {
		return cachedUser, nil
	}

	user, err := u.Repo.GetByID(ctx, userID)
	if err != nil {
		return nil, MapRepoError(err)
	}

	if err = u.Cache.Set(ctx, user); err != nil {
		return nil, MapRepoError(err)
	}

	user.SanitizePassword()
	return user, nil
}

// Find users by username
func (u *User) FindByUsername(ctx context.Context, name string, query *util.PaginationQuery) (*model.UsersList, error) {
	return u.Repo.FindByName(ctx, name, query)
}

// Get users with pagination
func (u *User) GetUsers(ctx context.Context, pq *util.PaginationQuery) (*model.UsersList, error) {
	return u.Repo.GetUsers(ctx, pq)
}
