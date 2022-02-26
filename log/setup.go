/*================================================================
*   Copyright (C) 2022. All rights reserved.
*
*   file : setup.go
*   coder: zemanzeng
*   date : 2022-01-16 18:04:28
*   desc : logger setup
*
================================================================*/

package log

import (
	"errors"
	"os"
)

const (
	StdLogName     = "__std"
	DefaultLogName = "__default"
)

var (
	writers map[string]Writer
	loggers map[string]Logger

	defaultLogger Logger
)

func init() {
	writers = make(map[string]Writer)
	loggers = make(map[string]Logger)

	RegisterWithWriter(DefaultLogName, LevelStrings[LevelDebug], NewWriter(os.Stdout))
}

func DefaultLogger() Logger {
	return Get(DefaultLogName)
}

func Register(name string, logger Logger) {
	loggers[name] = logger

	if name == DefaultLogName {
		defaultLogger = logger
	}
}

func RegisterWithFileName(name string, level string, fileName string) error {
	return RegisterWithLogConfig(name, NewLogConfig(level, fileName))
}

func RegisterWithLogConfig(name string, cfg *LogConfig) error {
	if cfg == nil || cfg.Base == nil || cfg.Writer == nil {
		return errors.New("invalid logger:" + name + " config")
	}

	var writer Writer
	if !cfg.Base.Console && len(cfg.Writer.FileName) > 0 {
		writer = NewSyncWriter(cfg.Writer)
		if cfg.Async.AsyncWrite {
			writer = NewAsyncWriter(writer, cfg.Async)
		}
	} else {
		writer = NewWriter(os.Stdout)
	}

	if _, err := writer.Write([]byte("init logger:" + name + " succ\n")); err != nil {
		Errorf("init logger:%s cfg:%+v error:%v", name, cfg, err)
		return err
	}

	RegisterWithWriter(name, cfg.Base.Level, writer)
	return nil
}

func RegisterWithWriter(name string, level string, writer Writer) {
	writers[name] = writer
	adapter := NewAdapter(level, writer, 0)
	loggers[name] = adapter

	if name == DefaultLogName {
		defaultLogger = adapter.WithDeep(1)
	}
}

func Get(name string) Logger {
	if logger, exist := loggers[name]; exist {
		return logger
	}
	return loggers[DefaultLogName]
}

func GetWriter(name string) Writer {
	if writer, exist := writers[name]; exist {
		return writer
	}
	return writers[DefaultLogName]
}

func Sync() {
	for name, logger := range loggers {
		if err := logger.Sync(); err != nil {
			Errorf("logger:%v sync error:", name, err)
		}
	}
}
