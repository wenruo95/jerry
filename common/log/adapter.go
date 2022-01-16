/*================================================================
*   Copyright (C) 2022. All rights reserved.
*
*   file : adapter.go
*   coder: zemanzeng
*   date : 2022-01-16 21:09:08
*   desc :
*
================================================================*/

package log

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	traceTag = "TRACE"
	debugTag = "DEBUG"
	infoTag  = "INFO"
	warnTag  = "WARN"
	errorTag = "ERROR"
	fatalTag = "FATAL"
)

type Adapter struct {
	lv     Level
	writer Writer
	deep   int
	fields []string
	mu     sync.RWMutex
}

func NewAdapter(lv string, writer Writer, deep int) *Adapter {
	return &Adapter{
		lv:     LevelNames[lv],
		writer: writer,
		deep:   deep,
		fields: make([]string, 0),
	}
}

func (ad *Adapter) Trace(v ...interface{}) {
	if ad.lv < LevelTrace {
		return
	}
	ad.Output(ad.deep+2, traceTag, fmt.Sprint(v...))
}

func (ad *Adapter) Tracef(format string, v ...interface{}) {
	if ad.lv < LevelTrace {
		return
	}
	ad.Output(ad.deep+2, traceTag, fmt.Sprintf(format, v...))
}

func (ad *Adapter) Debug(v ...interface{}) {
	if ad.lv < LevelDebug {
		return
	}
	ad.Output(ad.deep+2, debugTag, fmt.Sprint(v...))
}

func (ad *Adapter) Debugf(format string, v ...interface{}) {
	if ad.lv < LevelDebug {
		return
	}
	ad.Output(ad.deep+2, debugTag, fmt.Sprintf(format, v...))
}

func (ad *Adapter) Info(v ...interface{}) {
	if ad.lv < LevelInfo {
		return
	}
	ad.Output(ad.deep+2, infoTag, fmt.Sprint(v...))
}

func (ad *Adapter) Infof(format string, v ...interface{}) {
	if ad.lv < LevelInfo {
		return
	}
	ad.Output(ad.deep+2, infoTag, fmt.Sprintf(format, v...))
}

func (ad *Adapter) Warn(v ...interface{}) {
	if ad.lv < LevelWarn {
		return
	}
	ad.Output(ad.deep+2, warnTag, fmt.Sprint(v...))
}

func (ad *Adapter) Warnf(format string, v ...interface{}) {
	if ad.lv < LevelWarn {
		return
	}
	ad.Output(ad.deep+2, warnTag, fmt.Sprintf(format, v...))
}

func (ad *Adapter) Error(v ...interface{}) {
	if ad.lv < LevelError {
		return
	}
	ad.Output(ad.deep+2, errorTag, fmt.Sprint(v...))
}

func (ad *Adapter) Errorf(format string, v ...interface{}) {
	if ad.lv < LevelError {
		return
	}
	ad.Output(ad.deep+2, errorTag, fmt.Sprintf(format, v...))
}

func (ad *Adapter) Fatal(v ...interface{}) {
	if ad.lv < LevelFatal {
		return
	}
	ad.Output(ad.deep+2, fatalTag, fmt.Sprint(v...))
	os.Exit(1)
}

func (ad *Adapter) Fatalf(format string, v ...interface{}) {
	if ad.lv < LevelFatal {
		return
	}
	ad.Output(ad.deep+2, fatalTag, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func (ad *Adapter) Sync() error {
	return ad.writer.Sync()
}

func (ad *Adapter) SetLevel(level string) {
	lv, exist := LevelNames[strings.ToLower(level)]
	if !exist {
		return
	}

	ad.mu.Lock()
	ad.lv = lv
	ad.mu.Unlock()
}

func (ad *Adapter) GetLevel() string {
	ad.mu.RLock()
	level := ad.lv
	ad.mu.RUnlock()

	return level.String()
}

func (ad *Adapter) WithFields(fields ...string) Logger {
	ad.mu.RLock()
	defer ad.mu.RUnlock()
	adapter := NewAdapter(ad.lv.String(), ad.writer, ad.deep)
	adapter.fields = append(ad.fields, fields[0:len(fields)/2*2]...)
	return adapter
}

func (ad *Adapter) Output(calldepth int, tag, s string) error {
	now := time.Now() // get this early.

	_, file, line, ok := runtime.Caller(calldepth)
	if !ok {
		file = "???"
		line = 0
	}

	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	file = short

	year, month, day := now.Date()
	hour, min, sec := now.Clock()
	micro := now.Nanosecond() / 1e3

	suffix := ""
	for i := 0; i < len(ad.fields); i = i + 2 {
		pairs := "\"" + ad.fields[i] + "\":\"" + ad.fields[i+1] + "\""
		if i == 0 {
			suffix = "{" + pairs
		} else {
			suffix = suffix + "," + pairs
		}
	}

	if len(tag) > 0 {
		tag = "[" + tag + "] "
	}
	_, err := fmt.Fprintf(ad.writer, "%d-%d-%d %d:%d:%d.%d %s %s:%d %s\t%s",
		year, month, day, hour, min, sec, micro, tag, file, line, s, suffix)
	return err
}
