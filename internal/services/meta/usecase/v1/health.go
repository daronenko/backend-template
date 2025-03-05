package usecase

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

var (
	ErrPostgresNotReachable = errors.New("postgres not reachable")
	ErrRedisNotReachable    = errors.New("redis not reachable")
)

type Health struct {
	fx.In

	DB    *sqlx.DB
	Redis *redis.Client
}

// Needs to provide to fx
func NewHealth(u Health) HealthUsecase {
	return &u
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
