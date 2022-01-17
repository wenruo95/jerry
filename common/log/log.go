/*================================================================
*   Copyright (C) 2022. All rights reserved.
*
*   file : log.go
*   coder: zemanzeng
*   date : 2022-01-16 21:03:45
*   desc :
*
================================================================*/

package log

import (
	"context"

	"github.com/wenruo95/jerry/common/codec"
)

func Trace(format string) {
	defaultLogger.Trace(format)
}

func Tracef(format string, v ...interface{}) {
	defaultLogger.Tracef(format, v...)
}

func Debug(format string) {
	defaultLogger.Debug(format)
}

func Debugf(format string, v ...interface{}) {
	defaultLogger.Debugf(format, v...)
}

func Info(format string) {
	defaultLogger.Info(format)
}

func Infof(format string, v ...interface{}) {
	defaultLogger.Infof(format, v...)
}

func Warn(format string) {
	defaultLogger.Debug(format)
}

func Wranf(format string, v ...interface{}) {
	defaultLogger.Debugf(format, v...)
}

func Error(format string) {
	defaultLogger.Error(format)
}

func Errorf(format string, v ...interface{}) {
	defaultLogger.Errorf(format, v...)
}

func Fatal(format string) {
	defaultLogger.Fatal(format)
}

func Fatalf(format string, v ...interface{}) {
	defaultLogger.Fatalf(format, v...)
}

func WithFields(fields ...string) Logger {
	return Get(DefaultTag).WithFields(fields...)
}

func WithContextFields(ctx context.Context, fields ...string) context.Context {
	newCtx, msg := codec.EnsureMsg(ctx)
	logger, ok := msg.Logger().(Logger)
	if ok {
		logger = logger.WithFields(fields...)
	} else {
		logger = defaultLogger.WithFields(fields...)
	}
	msg.SetLogger(logger)
	return newCtx
}

func TraceContext(ctx context.Context, v ...interface{}) {
	switch logger := codec.Message(ctx).Logger().(type) {
	case Logger:
		logger.Trace(v...)
	default:
		defaultLogger.Trace(v...)
	}
}

func TraceContextf(ctx context.Context, format string, v ...interface{}) {
	switch logger := codec.Message(ctx).Logger().(type) {
	case Logger:
		logger.Tracef(format, v...)
	default:
		defaultLogger.Tracef(format, v...)
	}
}

func DebugContext(ctx context.Context, v ...interface{}) {
	switch logger := codec.Message(ctx).Logger().(type) {
	case Logger:
		logger.Debug(v...)
	default:
		defaultLogger.Debug(v...)
	}
}

func DebugContextf(ctx context.Context, format string, v ...interface{}) {
	switch logger := codec.Message(ctx).Logger().(type) {
	case Logger:
		logger.Debugf(format, v...)
	default:
		defaultLogger.Debugf(format, v...)
	}
}

func InfoContext(ctx context.Context, v ...interface{}) {
	switch logger := codec.Message(ctx).Logger().(type) {
	case Logger:
		logger.Info(v...)
	default:
		defaultLogger.Info(v...)
	}
}

func InfoContextf(ctx context.Context, format string, v ...interface{}) {
	switch logger := codec.Message(ctx).Logger().(type) {
	case Logger:
		logger.Infof(format, v...)
	default:
		defaultLogger.Infof(format, v...)
	}
}

func WarnContext(ctx context.Context, v ...interface{}) {
	switch logger := codec.Message(ctx).Logger().(type) {
	case Logger:
		logger.Warn(v...)
	default:
		defaultLogger.Warn(v...)
	}
}

func WarnContextf(ctx context.Context, format string, v ...interface{}) {
	switch logger := codec.Message(ctx).Logger().(type) {
	case Logger:
		logger.Warnf(format, v...)
	default:
		defaultLogger.Warnf(format, v...)
	}
}

func ErrorContext(ctx context.Context, v ...interface{}) {
	switch logger := codec.Message(ctx).Logger().(type) {
	case Logger:
		logger.Error(v...)
	default:
		defaultLogger.Error(v...)
	}
}

func ErrorContextf(ctx context.Context, format string, v ...interface{}) {
	switch logger := codec.Message(ctx).Logger().(type) {
	case Logger:
		logger.Errorf(format, v...)
	default:
		defaultLogger.Errorf(format, v...)
	}
}

func FatalContext(ctx context.Context, v ...interface{}) {
	switch logger := codec.Message(ctx).Logger().(type) {
	case Logger:
		logger.Fatal(v...)
	default:
		defaultLogger.Fatal(v...)
	}
}

func FatalContextf(ctx context.Context, format string, v ...interface{}) {
	switch logger := codec.Message(ctx).Logger().(type) {
	case Logger:
		logger.Fatalf(format, v...)
	default:
		defaultLogger.Fatalf(format, v...)
	}
}
