package config

import (
	"flag"
	"strings"

	"emperror.dev/errors"
	"github.com/daronenko/backend-template/internal/logger"
	"github.com/spf13/viper"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "", "path to service config file")
}

type Config struct {
	Service  ServiceConfig  `mapstructure:"service"`
	Postgres PostgresConfig `mapstructure:"postgres"`
}

type ServiceConfig struct {
	Name   string         `mapstructure:"name"`
	Logger *logger.Config `mapstructure:"logger"`
}

type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	SSLMode  bool   `mapstructure:"sslMode"`
}

func New() (*Config, error) {
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "viper.ReadInConfig")
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	for _, key := range viper.AllKeys() {
		value := viper.Get(key)
		viper.Set(key, value)
	}

	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, errors.Wrap(err, "viper.Unmarshal")
	}

	return cfg, nil
}
