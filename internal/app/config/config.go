package config

import "github.com/daronenko/backend-template/internal/app/ctx"

type Config struct {
	Service  ServiceSpec  `mapstructure:"service"`
	Postgres PostgresSpec `mapstructure:"postgres"`
	Redis    RedisSpec    `mapstructure:"redis"`

	AppContext ctx.Ctx
}
