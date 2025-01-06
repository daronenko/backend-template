package zl

import (
	"context"
	"fmt"

	"github.com/daronenko/backend-template/internal/logger"

	"github.com/rs/zerolog"
)

func (l *zerologLogger) WithContext(ctx context.Context) *zerologLogger {
	fields := zerolog.Ctx(ctx)
	return &zerologLogger{
		logger: l.logger.With().Fields(fields).Logger(),
	}
}

func (l *zerologLogger) WithFields(fields logger.Fields) *zerologLogger {
	return &zerologLogger{
		logger: l.logger.With().Fields(fields).Logger(),
	}
}

func (l *zerologLogger) Trace(args ...interface{}) {
	l.logger.Trace().Msg(fmt.Sprint(args...))
}

func (l *zerologLogger) Tracef(template string, args ...interface{}) {
	l.logger.Trace().Msgf(template, args...)
}

func (l *zerologLogger) Tracew(msg string, fields logger.Fields) {
	l.logger.Trace().Fields(fields).Msg(msg)
}

func (l *zerologLogger) Debug(args ...interface{}) {
	l.logger.Debug().Msg(fmt.Sprint(args...))
}

func (l *zerologLogger) Debugf(template string, args ...interface{}) {
	l.logger.Debug().Msgf(template, args...)
}

func (l *zerologLogger) Debugw(msg string, fields logger.Fields) {
	l.logger.Debug().Fields(fields).Msg(msg)
}

func (l *zerologLogger) Info(args ...interface{}) {
	l.logger.Info().Msg(fmt.Sprint(args...))
}

func (l *zerologLogger) Infof(template string, args ...interface{}) {
	l.logger.Info().Msgf(template, args...)
}

func (l *zerologLogger) Infow(msg string, fields logger.Fields) {
	l.logger.Info().Fields(fields).Msg(msg)
}

func (l *zerologLogger) Warn(args ...interface{}) {
	l.logger.Warn().Msg(fmt.Sprint(args...))
}

func (l *zerologLogger) Warnf(template string, args ...interface{}) {
	l.logger.Warn().Msgf(template, args...)
}

func (l *zerologLogger) Warnw(msg string, fields logger.Fields) {
	l.logger.Warn().Fields(fields).Msg(msg)
}

func (l *zerologLogger) WarnErr(msg string, err error) {
	l.logger.Warn().Stack().Err(err).Msg(msg)
}

func (l *zerologLogger) Error(args ...interface{}) {
	l.logger.Error().Msg(fmt.Sprint(args...))
}

func (l *zerologLogger) Errorf(template string, args ...interface{}) {
	l.logger.Error().Msgf(template, args...)
}

func (l *zerologLogger) Errorw(msg string, fields logger.Fields) {
	l.logger.Error().Fields(fields).Msg(msg)
}

func (l *zerologLogger) ErrorErr(msg string, err error) {
	l.logger.Error().Stack().Err(err).Msg(msg)
}

func (l *zerologLogger) Panic(args ...interface{}) {
	l.logger.Panic().Msg(fmt.Sprint(args...))
}

func (l *zerologLogger) Panicf(template string, args ...interface{}) {
	l.logger.Panic().Msgf(template, args...)
}

func (l *zerologLogger) Panicw(msg string, fields logger.Fields) {
	l.logger.Panic().Fields(fields).Msg(msg)
}

func (l *zerologLogger) PanicErr(msg string, err error) {
	l.logger.Panic().Stack().Err(err).Msg(msg)
}

func (l *zerologLogger) Fatal(args ...interface{}) {
	l.logger.Fatal().Msg(fmt.Sprint(args...))
}

func (l *zerologLogger) Fatalf(template string, args ...interface{}) {
	l.logger.Fatal().Msgf(template, args...)
}

func (l *zerologLogger) Fatalw(msg string, fields logger.Fields) {
	l.logger.Fatal().Fields(fields).Msg(msg)
}

func (l *zerologLogger) FatalErr(msg string, err error) {
	l.logger.Fatal().Stack().Err(err).Msg(msg)
}
