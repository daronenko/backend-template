package usecase

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

var (
	ErrPostgresNotReachable = errors.New("postgres not reachable")
	ErrRedisNotReachable    = errors.New("redis not reachable")
)

type Health struct {
	DB    *sqlx.DB
	Redis *redis.Client
}

func NewHealth(db *sqlx.DB, redis *redis.Client) *Health {
	return &Health{db, redis}
}

func (u *Health) Ping(ctx context.Context) error {
	if err := u.DB.PingContext(ctx); err != nil {
		return errors.Wrap(ErrPostgresNotReachable, err.Error())
	}

	if err := u.Redis.Ping(ctx).Err(); err != nil {
		return errors.Wrap(ErrRedisNotReachable, err.Error())
	}

	return nil
}
