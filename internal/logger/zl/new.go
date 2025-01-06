package zl

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/daronenko/backend-template/internal/logger"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/natefinch/lumberjack.v2"
)

type zerologLogger struct {
	logger zerolog.Logger
}

func New(cfg *logger.Config) *zerologLogger {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.DurationFieldUnit = time.Nanosecond
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	level, err := zerolog.ParseLevel(strings.ToLower(cfg.Level))
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	zl := zerolog.New(getWriter(cfg)).
		Level(level).
		With().
		Timestamp().
		Caller().
		Logger()

	return &zerologLogger{zl}
}

func getWriter(cfg *logger.Config) io.Writer {
	var writer io.Writer
	if cfg.DevMode {
		writer = os.Stdout
	} else {
		writer = zerolog.MultiLevelWriter(
			&lumberjack.Logger{
				Filename:   cfg.LogsPath,
				MaxSize:    cfg.MaxSize,
				MaxBackups: cfg.MaxBackups,
				MaxAge:     cfg.MaxAge,
				Compress:   cfg.Compress,
			},
			zerolog.ConsoleWriter{
				Out:        os.Stdout,
				TimeFormat: time.RFC3339Nano,
			},
		)
	}

	return writer
}
