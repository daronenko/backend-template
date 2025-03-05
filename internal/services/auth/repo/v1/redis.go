package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"emperror.dev/errors"
	"github.com/daronenko/backend-template/internal/app/config"
	"github.com/daronenko/backend-template/internal/model/v1"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

// Cache for users
type UserCache struct {
	fx.In

	Redis *redis.Client
	Conf  *config.Config
}

// Needs to provide to fx
func NewUserCache(c UserCache) Cache {
	return &c
}

// Cache user
func (r *UserCache) Set(ctx context.Context, user *model.User) error {
	userBytes, err := json.Marshal(user)
	if err != nil {
		return errors.Wrap(err, "repo.UserCache.Set.json.Marshal")
	}

	if err = r.Redis.Set(ctx, r.cacheKey(user.ID), userBytes, time.Second*time.Duration(r.Conf.App.Auth.User.Cache.Expire)).Err(); err != nil {
		return errors.Wrap(err, "repo.UserCache.Set.redis.Set")
	}

	return nil
}

// Get user by id from cache
func (r *UserCache) GetByID(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	userStr, err := r.Redis.Get(ctx, r.cacheKey(userID)).Result()
	if err != nil {
		return nil, errors.Wrap(err, "repo.UserCache.GetByID.redis.Get")
	}

	user := &model.User{}
	if err = json.Unmarshal([]byte(userStr), user); err != nil {
		return nil, errors.Wrap(err, "repo.UserCache.GetByID.json.Unmarshal")
	}

	return user, nil
}

// Delete user by id from cache
func (r *UserCache) Delete(ctx context.Context, userID uuid.UUID) error {
	if err := r.Redis.Del(ctx, r.cacheKey(userID)).Err(); err != nil {
		return errors.Wrap(err, "repo.UserCache.Delete.redis.Del")
	}

	return nil
}

func (r *UserCache) cacheKey(userID uuid.UUID) string {
	return fmt.Sprintf("%s:%s", r.Conf.App.Auth.User.Cache.Prefix, userID)
}
