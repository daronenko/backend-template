package logger

type Config struct {
	Level      string `mapstructure:"level"`
	DevMode    bool   `mapstructure:"devMode"`
	LogsPath   string `mapstructure:"logsPath"`
	MaxSize    int    `mapstructure:"maxSize"`
	MaxBackups int    `mapstructure:"maxBackups"`
	MaxAge     int    `mapstructure:"maxAge"`
	Compress   bool   `mapstructure:"compress"`
}
