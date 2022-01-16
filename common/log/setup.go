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

const DefaultTag = "default"

var (
	writers map[string]Writer
	loggers map[string]Logger

	StdLogger     Logger
	DefaultLogger Logger
)

func init() {
	writers = make(map[string]Writer)
	loggers = make(map[string]Logger)

	StdLogger = NewAdapter(LevelStrings[LevelDebug], NewWriter(os.Stdout), 0)
	DefaultLogger = NewAdapter(LevelStrings[LevelDebug], NewWriter(os.Stdout), 1)
}

func Register(tag string, level string, fileName string) {
	RegisterWithLogConfig(tag, NewLogConfig(level, fileName))
}

func RegisterWithLogConfig(tag string, cfg *LogConfig) {
	var writer Writer = NewSyncWriter(cfg.Writer)
	if cfg.Async != nil && cfg.Async.AsyncWrite {
		writer = NewAsyncWriter(writer, cfg.Async)
	}
	RegisterWithWriter(tag, cfg.Base.Level, writer)
}

func RegisterWithWriter(tag string, level string, writer Writer) {
	writers[tag] = writer
	loggers[tag] = NewAdapter(level, writer, 0)

	if tag == DefaultTag {
		DefaultLogger = NewAdapter(level, writer, 1)
	}
}

func Get(tag string) Logger {
	if logger, exist := loggers[tag]; exist {
		return logger
	}
	return loggers[DefaultTag]
}

func Setup(cfg *LogConfig) {
	if cfg == nil {
		return
	}
	RegisterWithLogConfig(DefaultTag, cfg)
}
