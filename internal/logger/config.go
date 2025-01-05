package logger

type Config struct {
	Level   string `mapstructure:"level"`
	DevMode bool   `mapstructure:"devMode"`
	Encoder string `mapstructure:"encoder"`
}
