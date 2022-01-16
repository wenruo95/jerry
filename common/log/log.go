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

import "context"

func Trace(format string) {
	DefaultLogger.Trace(format)
}

func Tracef(format string, v ...interface{}) {
	DefaultLogger.Tracef(format, v...)
}

func Debug(format string) {
	DefaultLogger.Debug(format)
}

func Debugf(format string, v ...interface{}) {
	DefaultLogger.Debugf(format, v...)
}

func Info(format string) {
	DefaultLogger.Info(format)
}

func Infof(format string, v ...interface{}) {
	DefaultLogger.Infof(format, v...)
}

func Warn(format string) {
	DefaultLogger.Debug(format)
}

func Wranf(format string, v ...interface{}) {
	DefaultLogger.Debugf(format, v...)
}

func Error(format string) {
	DefaultLogger.Error(format)
}

func Errorf(format string, v ...interface{}) {
	DefaultLogger.Errorf(format, v...)
}

func Fatal(format string) {
	DefaultLogger.Fatal(format)
}

func Fatalf(format string, v ...interface{}) {
	DefaultLogger.Fatalf(format, v...)
}

func WithContextFields(ctx context.Context, fields ...string) Logger {
	return DefaultLogger.WithFields(fields...)
}

// TODO
