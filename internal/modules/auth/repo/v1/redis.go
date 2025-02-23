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
)

// Cache for users
type UserCache struct {
	redis *redis.Client
	conf  *config.Config
}

func NewUserCache(redis *redis.Client, conf *config.Config) Cache {
	return &UserCache{redis, conf}
}

// Cache user
func (r *UserCache) Set(ctx context.Context, user *model.User) error {
	userBytes, err := json.Marshal(user)
	if err != nil {
		return errors.Wrap(err, "repo.UserCache.Set.json.Marshal")
	}

	if err = r.redis.Set(ctx, r.cacheKey(user.ID), userBytes, time.Second*time.Duration(r.conf.App.Auth.User.Cache.Expire)).Err(); err != nil {
		return errors.Wrap(err, "repo.UserCache.Set.redis.Set")
	}

	return nil
}

// Get user by id from cache
func (r *UserCache) GetByID(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	userStr, err := r.redis.Get(ctx, r.cacheKey(userID)).Result()
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
	if err := r.redis.Del(ctx, r.cacheKey(userID)).Err(); err != nil {
		return errors.Wrap(err, "repo.UserCache.Delete.redis.Del")
	}

	return nil
}

func (r *UserCache) cacheKey(userID uuid.UUID) string {
	return fmt.Sprintf("%s:%s", r.conf.App.Auth.User.Cache.Prefix, userID)
}
