package zl

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/daronenko/backend-template/internal/config"
	"github.com/daronenko/backend-template/internal/constants"
	"github.com/daronenko/backend-template/internal/logger"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/natefinch/lumberjack.v2"
)

type zerologLogger struct {
	logger zerolog.Logger
}

func New(cfg *logger.Config) zerologLogger {
	level, err := zerolog.ParseLevel(strings.ToLower(cfg.Level))
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	var output io.Writer
	if cfg.DevMode {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
	} else {
		output = &lumberjack.Logger{
			Filename:   "/var/log/dump.log",
			MaxSize:    10,
			MaxBackups: 5,
			MaxAge:     30,
			Compress:   true,
		}
	}

	return zerologLogger{
		zerolog.New(output).
			Level(level).
			With().
			Timestamp().
			Str(constants.Revision, config.Revision).
			Str(constants.Version, config.Version).
			Logger(),
	}
}
