package config

import (
	"flag"
	"fmt"
	"strings"

	"github.com/daronenko/backend-template/internal/app/ctx"
	"github.com/spf13/viper"
)

func New(ctx ctx.Ctx) (*Config, error) {
	flag.Parse()

	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	for _, key := range viper.AllKeys() {
		value := viper.Get(key)
		viper.Set(key, value)
	}

	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return cfg, nil
}
