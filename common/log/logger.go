/*================================================================
*   Copyright (C) 2022. All rights reserved.
*
*   file : logger.go
*   coder: zemanzeng
*   date : 2022-01-16 11:46:24
*   desc : logger
*
================================================================*/

package log

type Level int

const (
	LevelNil Level = iota
	LevelTrace
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

func (lv *Level) String() string {
	return LevelStrings[*lv]
}

var LevelStrings map[Level]string = map[Level]string{
	LevelTrace: "trace",
	LevelDebug: "debug",
	LevelInfo:  "info",
	LevelWarn:  "warn",
	LevelError: "error",
	LevelFatal: "fatal",
}

var LevelNames map[string]Level = map[string]Level{
	"trace": LevelTrace,
	"debug": LevelDebug,
	"info":  LevelInfo,
	"warn":  LevelWarn,
	"error": LevelError,
	"fatal": LevelFatal,
}

type Logger interface {
	Trace(v ...interface{})
	Tracef(format string, v ...interface{})
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})

	Sync() error

	SetLevel(level string)
	GetLevel() string

	WithFields(field ...string) Logger
}

type Writer interface {
	Write(p []byte) (n int, err error)
	Sync() error
	Close() error
}
