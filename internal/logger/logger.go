package logger

import (
	"time"
)

type Fields map[string]interface{}

type Logger interface {
	InitLogger()
	Sync() error
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Debugw(msg string, fields Fields)
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Infow(msg string, fields Fields)
	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
	WarnMsg(msg string, err error)
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	Errorw(msg string, fields Fields)
	Err(msg string, err error)
	DPanic(args ...interface{})
	DPanicf(template string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
	Printf(template string, args ...interface{})
	WithName(name string)
	HttpMiddlewareAccessLogger(
		method string,
		uri string,
		status int,
		size int64,
		time time.Duration,
	)
	GrpcMiddlewareAccessLogger(
		method string,
		time time.Duration,
		metaData map[string][]string,
		err error,
	)
	GrpcClientInterceptorLogger(
		method string,
		req interface{},
		reply interface{},
		time time.Duration,
		metaData map[string][]string,
		err error,
	)
}
