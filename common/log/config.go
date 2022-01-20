/*================================================================
*   Copyright (C) 2022. All rights reserved.
*
*   file : config.go
*   coder: zemanzeng
*   date : 2022-01-16 17:46:55
*   desc : log config
*
================================================================*/

package log

// LogConfig 日志配置
type LogConfig struct {
	Base   *BaseConfig   `json:"base" yaml:"base"`     // 基础配置
	Writer *WriterConfig `json:"writer" yaml:"writer"` // 写配置
	Async  *AsyncConfig  `json:"async" yaml:"async"`   // 异步配置
}

// BaseConfig 日志基础配置
type BaseConfig struct {
	Level   string `json:"level" yaml:"level"`     // 日志级别 trace debug info warn error fatal
	Console bool   `json:"console" yaml:"console"` // 是否输出到控制台
}

// WriterConfig 写配置
type WriterConfig struct {
	FileName   string `json:"file_name,filename" yaml:"file_name,filename"`         // 日志文件存放目录 默认:""
	MaxSize    int    `json:"max_size,maxsize" yaml:"max_size,maxsize"`             // 文件大小限制 单位:MB 默认:10M
	MaxBackups int    `json:"max_backups,maxbackups" yaml:"max_backups,maxbackups"` // 最大保留日志文件数量 默认:3个
	MaxAge     int    `json:"max_age,maxage" yaml:"max_age,maxage"`                 // 日志文件保留天数 默认:3天
	Compress   bool   `json:"compress" yaml:"compress"`                             // 是否压缩处理 默认:false
	LocalTime  bool   `json:"local_time,localtime" yaml:"local_time,localtime"`     // 是否用本地时间 默认true false:UTC时间
}

// AsynerConfig 异步配置
type AsyncConfig struct {
	AsyncWrite bool  `json:"async_write,asyncwrite" yaml:"async_write,asyncwrite"` // 是否异步 默认:否
	DropLog    bool  `json:"drop_log,droplog" yaml:"drop_log,droplog"`             // 有阻塞时是否丢弃日志 默认:否
	QueueSize  int   `json:"queue_size,queuesize" yaml:"queue_size,queuesize"`     // 队列大小 默认:10000
	WriteSize  int   `json:"write_size,writesize" yaml:"write_size,writesize"`     // 多大刷盘一次 单位byte 默认:4k
	IntervalMs int64 `json:"interval_ms,intervalms" yaml:"interval_ms,intervalms"` // 刷盘间隔 单位:ms 默认:100ms
}

func NewLogConfig(level string, fileName string) *LogConfig {
	cfg := &LogConfig{
		Base: &BaseConfig{
			Level:   LevelStrings[LevelDebug],
			Console: true,
		},
		Writer: &WriterConfig{
			FileName:   fileName,
			MaxSize:    10,
			MaxBackups: 3,
			MaxAge:     3,
			Compress:   false,
			LocalTime:  true,
		},
		Async: &AsyncConfig{
			AsyncWrite: false,
			DropLog:    false,
			QueueSize:  10000,
			IntervalMs: 100,
		},
	}

	if len(fileName) > 0 {
		cfg.Base.Console = false
	}

	return cfg
}
