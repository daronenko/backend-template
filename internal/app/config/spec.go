package config

import (
	"time"

	"github.com/daronenko/backend-template/pkg/logger"
)

type AppSpec struct {
	Address string        `mapstructure:"address"`
	Admin   AdminSpec     `mapstructure:"admin"`
	Logger  logger.Config `mapstructure:"logger"`
	Tracing TracingSpec   `mapstructure:"tracing"`
	Auth    AuthSpec      `mapstructure:"auth"`
}

type AdminSpec struct {
	Key string `mapstructure:"key"`
}

type TracingSpec struct {
	Enabled bool `mapstructure:"enabled" split_words:"true"`

	// TracingExporters to indicate which exporters to use for tracing.
	// Valid values are: jaeger, otlp, stdout (for debug).
	Exporters []string `mapstructure:"exporters" split_words:"true" default:"jaeger"`

	// TracingSampleRate to indicate the sampling rate for tracing.
	// Valid values are: 0.0 (disabled), 1.0 (all traces), or a value between 0.0 and 1.0 (sampling rate).
	SampleRate float64 `mapstructure:"sampleRate" split_words:"true" default:"1.0"`
}

type AuthSpec struct {
	Jwt     JwtSpec     `mapstructure:"jwt"`
	User    UserSpec    `mapstructure:"user"`
	Session SessionSpec `mapstructure:"session"`
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
