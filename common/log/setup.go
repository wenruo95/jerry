/*================================================================
*   Copyright (C) 2022. All rights reserved.
*
*   file : setup.go
*   coder: zemanzeng
*   date : 2022-01-16 18:04:28
*   desc :
*
================================================================*/

package log

import "os"

const (
	StdTag     = "__std"
	DefaultTag = "__default"
)

var (
	writers map[string]Writer
	loggers map[string]Logger

	defaultLogger Logger
)

func init() {
	writers = make(map[string]Writer)
	loggers = make(map[string]Logger)

	RegisterWithWriter(StdTag, LevelStrings[LevelDebug], NewWriter(os.Stdout))
	RegisterWithWriter(DefaultTag, LevelStrings[LevelDebug], NewWriter(os.Stdout))
}

func Register(tag string, level string, fileName string) {
	RegisterWithLogConfig(tag, NewLogConfig(level, fileName))
}

func RegisterWithLogConfig(tag string, cfg *LogConfig) {
	var writer Writer
	if !cfg.Base.Console && len(cfg.Writer.FileName) > 0 {
		writer = NewSyncWriter(cfg.Writer)
		if cfg.Async != nil && cfg.Async.AsyncWrite {
			writer = NewAsyncWriter(writer, cfg.Async)
		}
	} else {
		writer = NewWriter(os.Stdout)
	}
	RegisterWithWriter(tag, cfg.Base.Level, writer)
}

func RegisterWithWriter(tag string, level string, writer Writer) {
	writers[tag] = writer
	adapter := NewAdapter(level, writer, 0)
	loggers[tag] = adapter

	if tag == DefaultTag {
		defaultLogger = adapter.WithDeep(1)
	}
}

func Get(tag string) Logger {
	if logger, exist := loggers[tag]; exist {
		return logger
	}
	return loggers[DefaultTag]
}

func GetWriter(tag string) Writer {
	if writer, exist := writers[tag]; exist {
		return writer
	}
	return writers[DefaultTag]
}

func Setup(cfg *LogConfig) {
	if cfg == nil {
		return
	}
	RegisterWithLogConfig(DefaultTag, cfg)
}
