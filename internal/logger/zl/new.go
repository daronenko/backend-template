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

func New(cfg *logger.Config) (*zerologLogger, error) {
	level, err := zerolog.ParseLevel(strings.ToLower(cfg.Level))
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	output, err := configureOutput(cfg)
	if err != nil {
		return nil, err
	}

	zl := zerolog.New(output).
		Level(level).
		With().
		Timestamp().
		Logger()

	return &zerologLogger{logger: zl}, nil
}

func configureOutput(cfg *logger.Config) (io.Writer, error) {
	if cfg.DevMode {
		return zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}, nil
	}

	return &lumberjack.Logger{
		Filename:   cfg.LogsPath,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}, nil
}
