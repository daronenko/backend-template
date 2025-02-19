package infra

import (
	"context"
	"fmt"
	"time"

	"github.com/daronenko/backend-template/internal/app/config"
	"github.com/redis/go-redis/v9"
)

func NewRedis(conf *config.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", conf.Redis.Host, conf.Redis.Port),
		Password:     conf.Redis.Password,
		DB:           conf.Redis.Database,
		MinIdleConns: conf.Redis.MinIdleConns,
		PoolSize:     conf.Redis.PoolSize,
		PoolTimeout:  time.Duration(conf.Redis.PoolTimeout) * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	ping := client.Ping(ctx)
	if ping.Err() != nil {
		return nil, ping.Err()
	}

	return client, nil
}
