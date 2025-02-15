package zerologger

import (
	"io"
	"log"
	"os"
	"strings"
	"time"

	. "github.com/daronenko/backend-template/pkg/logger/config"
	. "github.com/daronenko/backend-template/pkg/logger/contracts"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/natefinch/lumberjack.v2"
)

type zeroLogger struct {
	logger zerolog.Logger
}

func New(cfg *Config) Logger {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.DurationFieldUnit = time.Nanosecond
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	level, err := zerolog.ParseLevel(strings.ToLower(cfg.Level))
	if err != nil {
		log.Printf("error: failed to parse log level: %q\n", cfg.Level)
		os.Exit(1)
	}
	zerolog.SetGlobalLevel(level)

	l := zerolog.New(getWriter(cfg)).
		Level(level).
		With().
		Timestamp().
		Caller().
		Logger()

	return &zeroLogger{l}
}

func getWriter(cfg *Config) io.Writer {
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
