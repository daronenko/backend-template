package config

import "github.com/daronenko/backend-template/pkg/logger"

type ServiceSpec struct {
	Name     string        `mapstructure:"name"`
	AdminKey string        `mapstructure:"adminKey"`
	Logger   logger.Config `mapstructure:"logger"`
	Auth     AuthSpec      `mapstructure:"auth"`
}

type AuthSpec struct {
	JwtSecret string      `mapstructure:"jwtSecret"`
	User      UserSpec    `mapstructure:"user"`
	Session   SessionSpec `mapstructure:"session"`
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

type PostgresSpec struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	SSLMode  string `mapstructure:"sslMode"`
	Driver   string `mapstructure:"driver"`
}

type RedisSpec struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Password     string `mapstructure:"password"`
	Database     int    `mapstructure:"database"`
	MinIdleConns int    `mapstructure:"minIdleConns"`
	PoolSize     int    `mapstructure:"poolSize"`
	PoolTimeout  int    `mapstructure:"poolTimeout"`
}
