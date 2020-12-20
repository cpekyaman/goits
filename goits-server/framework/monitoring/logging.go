package monitoring

import (
	"go.uber.org/zap"
)

const (
	LF_CID        = "cid"
	LF_REQID      = "rqid"
	LF_Resource   = "res"
	LF_Operation  = "oper"
	LF_HttpStatus = "hsc"
	LF_Duration   = "durn"
	LF_Ms         = "ms"
	LF_Bytes      = "bytes"
	LF_SubSystem  = "sys"
	LF_Error      = "err"
	LF_ErrorType  = "ert"
	LF_ErrorCause = "eca"
	LF_Type       = "tp"
	LF_Query      = "q"
)

var logger *zap.Logger

func init() {
	logger, _ = zap.NewDevelopment(zap.AddCallerSkip(1))
}

// IntLogField returns a field with int value for adding to log context.
func IntLogField(key string, value int) zap.Field {
	return zap.Int(key, value)
}

// Int64LogField returns a field with int64 value for adding to log context.
func Int64LogField(key string, value int64) zap.Field {
	return zap.Int64(key, value)
}

// StrLogField returns a field with string value for adding to log context.
func StrLogField(key string, value string) zap.Field {
	return zap.String(key, value)
}

// ErrLogField returns a field with error as value for adding to log context.
func ErrLogField(err error) zap.Field {
	return zap.Error(err)
}

// RootLogger returns the vanilla logger for the application.
func RootLogger() Logger {
	return zapLogger{
		l: logger,
	}
}

// NewLogger creates a new logger with the given fields applied to log context.
func NewLogger(fields ...zap.Field) Logger {
	return zapLogger{
		l: logger.With(fields...),
	}
}

// NewLoggerWithMonitoring creates a new logger by using context parameters to setup logging context.
func NewLoggerWithMonitoring(mctx MonitoringContext) Logger {
	return zapLogger{
		l: logger.With(zap.String(LF_Resource, mctx.Resource()),
			zap.String(LF_Operation, mctx.Operation()),
			zap.String(LF_CID, mctx.CID()),
			zap.String(LF_REQID, mctx.REQID())),
	}
}

// Logger is a wrapper interface to reduce dependency on actual logger library to some extend.
type Logger interface {
	// WithStr returns a new logger with the given key-value added to context.
	WithStr(key string, value string) Logger

	// With returns a new logger with the given fields added to context.
	With(fields ...zap.Field) Logger

	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Fatal(msg string)
}

// zapLogger is the current logger implementation that uses uber's zap library.
type zapLogger struct {
	l *zap.Logger
}

func (this zapLogger) WithStr(key string, value string) Logger {
	return zapLogger{
		l: this.l.With(zap.String(key, value)),
	}
}

func (this zapLogger) With(fields ...zap.Field) Logger {
	return zapLogger{
		l: this.l.With(fields...),
	}
}

func (this zapLogger) Info(msg string) {
	this.l.Info(msg)
}

func (this zapLogger) Warn(msg string) {
	this.l.Warn(msg)
}

func (this zapLogger) Error(msg string) {
	this.l.Error(msg)
}

func (this zapLogger) Fatal(msg string) {
	this.l.Fatal(msg)
}
