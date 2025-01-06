package logger

type Fields map[string]interface{}

type Logger interface {
	WithName(name string)

	Trace(args ...interface{})
	Tracef(template string, args ...interface{})
	Tracew(msg string, fields Fields)

	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Debugw(msg string, fields Fields)

	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Infow(msg string, fields Fields)

	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
	Warnw(msg string, fields Fields)
	WarnErr(msg string, err error)

	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	Errorw(msg string, fields Fields)
	ErrorErr(msg string, err error)

	Panic(args ...interface{})
	Panicf(template string, args ...interface{})
	Panicw(msg string, fields Fields)
	PanicErr(msg string, err error)

	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
	Fatalw(msg string, fields Fields)
	FatalErr(msg string, err error)
}
