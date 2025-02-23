package config

import (
	"time"

	"github.com/daronenko/backend-template/pkg/logger"
)

type AppSpec struct {
	Address string        `mapstructure:"address"`
	Admin   AdminSpec     `mapstructure:"admin"`
	Logger  logger.Config `mapstructure:"logger"`
	Auth    AuthSpec      `mapstructure:"auth"`
}

type AdminSpec struct {
	Key string `mapstructure:"key"`
}

type AuthSpec struct {
	JwtSecret string      `mapstructure:"jwt"`
	User      UserSpec    `mapstructure:"user"`
	Session   SessionSpec `mapstructure:"session"`
}

type JwtSpec struct {
	Secret string `mapstructure:"secret"`
}

type UserSpec struct {
	Cache CacheSpec `mapstructure:"cache"`
}

type SessionSpec struct {
	Cache  CacheSpec  `mapstructure:"cache"`
	Cookie CookieSpec `mapstructure:"cookie"`
}

type CookieSpec struct {
	Name     string `mapstructure:"name"`
	Secure   bool   `mapstructure:"secure"`
	HTTPOnly bool   `mapstructure:"httpOnly"`
}

type CacheSpec struct {
	Prefix string `mapstructure:"prefix"`
	Expire int    `mapstructure:"expire"`
}

type DevOpsSpec struct {
	Address string `mapstructure:"address"`
}

type ServerSpec struct {
	ShutdownTimeout time.Duration `mapstructure:"shutdownTimeout"`
	TrustedProxies  []string      `mapstructure:"trustedProxies" split_words:"true"`
}

type PostgresSpec struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	SSLMode  string `mapstructure:"sslMode"`
	Driver   string `mapstructure:"driver"`

	MaxOpenConns    int           `mapstructure:"maxOpenConns"`
	ConnMaxLifetime time.Duration `mapstructure:"connMaxLifetime"`
	MaxIdleConns    int           `mapstructure:"maxIdleConns"`
	ConnMaxIdleTime time.Duration `mapstructure:"connMaxIdleTime"`
}

type RedisSpec struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Database int    `mapstructure:"database"`

	MinIdleConns int           `mapstructure:"minIdleConns"`
	PoolSize     int           `mapstructure:"poolSize"`
	PoolTimeout  time.Duration `mapstructure:"poolTimeout"`
}
