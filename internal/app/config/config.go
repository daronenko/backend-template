package config

import "github.com/daronenko/backend-template/internal/app/ctx"

type Config struct {
	DevMode  bool         `mapstructure:"devMode"`
	App      AppSpec      `mapstructure:"app"`
	DevOps   DevOpsSpec   `mapstructure:"devops"`
	Server   ServerSpec   `mapstructure:"server"`
	Postgres PostgresSpec `mapstructure:"postgres"`
	Redis    RedisSpec    `mapstructure:"redis"`

	AppContext ctx.Ctx
}
