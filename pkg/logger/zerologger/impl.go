package zerologger

import (
	"context"
	"fmt"

	. "github.com/daronenko/backend-template/pkg/logger"

	"github.com/rs/zerolog"
)

func (l *zeroLogger) WithContext(ctx context.Context) Logger {
	fields := zerolog.Ctx(ctx)
	return &zeroLogger{
		logger: l.logger.With().Fields(fields).Logger(),
	}
}

func (l *zeroLogger) WithFields(fields Fields) Logger {
	return &zeroLogger{
		logger: l.logger.With().Fields(fields).Logger(),
	}
}

func (l *zeroLogger) Trace(args ...interface{}) {
	l.logger.Trace().Msg(fmt.Sprint(args...))
}

func (l *zeroLogger) Tracef(template string, args ...interface{}) {
	l.logger.Trace().Msgf(template, args...)
}

func (l *zeroLogger) Tracew(msg string, fields Fields) {
	l.logger.Trace().Fields(fields).Msg(msg)
}

func (l *zeroLogger) Debug(args ...interface{}) {
	l.logger.Debug().Msg(fmt.Sprint(args...))
}

func (l *zeroLogger) Debugf(template string, args ...interface{}) {
	l.logger.Debug().Msgf(template, args...)
}

func (l *zeroLogger) Debugw(msg string, fields Fields) {
	l.logger.Debug().Fields(fields).Msg(msg)
}

func (l *zeroLogger) Info(args ...interface{}) {
	l.logger.Info().Msg(fmt.Sprint(args...))
}

func (l *zeroLogger) Infof(template string, args ...interface{}) {
	l.logger.Info().Msgf(template, args...)
}

func (l *zeroLogger) Infow(msg string, fields Fields) {
	l.logger.Info().Fields(fields).Msg(msg)
}

func (l *zeroLogger) Warn(args ...interface{}) {
	l.logger.Warn().Msg(fmt.Sprint(args...))
}

func (l *zeroLogger) Warnf(template string, args ...interface{}) {
	l.logger.Warn().Msgf(template, args...)
}

func (l *zeroLogger) Warnw(msg string, fields Fields) {
	l.logger.Warn().Fields(fields).Msg(msg)
}

func (l *zeroLogger) WarnErr(msg string, err error) {
	l.logger.Warn().Stack().Err(err).Msg(msg)
}

func (l *zeroLogger) Error(args ...interface{}) {
	l.logger.Error().Msg(fmt.Sprint(args...))
}

func (l *zeroLogger) Errorf(template string, args ...interface{}) {
	l.logger.Error().Msgf(template, args...)
}

func (l *zeroLogger) Errorw(msg string, fields Fields) {
	l.logger.Error().Fields(fields).Msg(msg)
}

func (l *zeroLogger) ErrorErr(msg string, err error) {
	l.logger.Error().Stack().Err(err).Msg(msg)
}

func (l *zeroLogger) Panic(args ...interface{}) {
	l.logger.Panic().Msg(fmt.Sprint(args...))
}

func (l *zeroLogger) Panicf(template string, args ...interface{}) {
	l.logger.Panic().Msgf(template, args...)
}

func (l *zeroLogger) Panicw(msg string, fields Fields) {
	l.logger.Panic().Fields(fields).Msg(msg)
}

func (l *zeroLogger) PanicErr(msg string, err error) {
	l.logger.Panic().Stack().Err(err).Msg(msg)
}

func (l *zeroLogger) Fatal(args ...interface{}) {
	l.logger.Fatal().Msg(fmt.Sprint(args...))
}

func (l *zeroLogger) Fatalf(template string, args ...interface{}) {
	l.logger.Fatal().Msgf(template, args...)
}

func (l *zeroLogger) Fatalw(msg string, fields Fields) {
	l.logger.Fatal().Fields(fields).Msg(msg)
}

func (l *zeroLogger) FatalErr(msg string, err error) {
	l.logger.Fatal().Stack().Err(err).Msg(msg)
}
