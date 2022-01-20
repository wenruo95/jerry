/*================================================================
*   Copyright (C) 2022. All rights reserved.
*
*   file : adapter.go
*   coder: zemanzeng
*   date : 2022-01-16 21:09:08
*   desc : logger implements
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
	prefix string
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
	if ad.lv > LevelTrace {
		return
	}
	ad.Output(ad.deep+2, traceTag, fmt.Sprint(v...))
}

func (ad *Adapter) Tracef(format string, v ...interface{}) {
	if ad.lv > LevelTrace {
		return
	}
	ad.Output(ad.deep+2, traceTag, fmt.Sprintf(format, v...))
}

func (ad *Adapter) Debug(v ...interface{}) {
	if ad.lv > LevelDebug {
		return
	}
	ad.Output(ad.deep+2, debugTag, fmt.Sprint(v...))
}

func (ad *Adapter) Debugf(format string, v ...interface{}) {
	if ad.lv > LevelDebug {
		return
	}
	ad.Output(ad.deep+2, debugTag, fmt.Sprintf(format, v...))
}

func (ad *Adapter) Info(v ...interface{}) {
	if ad.lv > LevelInfo {
		return
	}
	ad.Output(ad.deep+2, infoTag, fmt.Sprint(v...))
}

func (ad *Adapter) Infof(format string, v ...interface{}) {
	if ad.lv > LevelInfo {
		return
	}
	ad.Output(ad.deep+2, infoTag, fmt.Sprintf(format, v...))
}

func (ad *Adapter) Warn(v ...interface{}) {
	if ad.lv > LevelWarn {
		return
	}
	ad.Output(ad.deep+2, warnTag, fmt.Sprint(v...))
}

func (ad *Adapter) Warnf(format string, v ...interface{}) {
	if ad.lv > LevelWarn {
		return
	}
	ad.Output(ad.deep+2, warnTag, fmt.Sprintf(format, v...))
}

func (ad *Adapter) Error(v ...interface{}) {
	if ad.lv > LevelError {
		return
	}
	ad.Output(ad.deep+2, errorTag, fmt.Sprint(v...))
}

func (ad *Adapter) Errorf(format string, v ...interface{}) {
	if ad.lv > LevelError {
		return
	}
	ad.Output(ad.deep+2, errorTag, fmt.Sprintf(format, v...))
}

func (ad *Adapter) Fatal(v ...interface{}) {
	if ad.lv > LevelFatal {
		return
	}
	ad.Output(ad.deep+2, fatalTag, fmt.Sprint(v...))
	os.Exit(1)
}

func (ad *Adapter) Fatalf(format string, v ...interface{}) {
	if ad.lv > LevelFatal {
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

func (ad *Adapter) SetPrefix(prefix string) {
	ad.mu.Lock()
	ad.prefix = prefix
	ad.mu.Unlock()
}

func (ad *Adapter) WithDeep(deep int) Logger {
	adapter := ad.clone()
	adapter.deep = deep
	return adapter
}

func (ad *Adapter) WithFields(fields ...string) Logger {
	adapter := ad.clone()
	adapter.fields = append(adapter.fields, fields[0:len(fields)/2*2]...)
	return adapter
}

func (ad *Adapter) clone() *Adapter {
	ad.mu.RLock()
	defer ad.mu.RUnlock()

	adapter := NewAdapter(ad.lv.String(), ad.writer, ad.deep)
	adapter.fields = append(adapter.fields, ad.fields...)
	adapter.prefix = ad.prefix
	return adapter
}

func (ad *Adapter) Output(calldepth int, tag, s string) error {
	now := time.Now() // get this early.

	_, file, line, ok := runtime.Caller(calldepth)
	if !ok {
		file = "???"
		line = 0
	}

	var indexs []int
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			indexs = append(indexs, i)
			if len(indexs) >= 2 {
				break
			}
		}
	}
	file = file[indexs[len(indexs)-1]+1:]

	year, month, day := now.Date()
	hour, min, sec := now.Clock()
	millsec := now.Nanosecond() / 1e6

	ad.mu.RLock()
	var content string
	if len(ad.prefix) > 0 {
		content = ad.prefix + " " + s
	} else {
		content = content + s
	}

	var pairs string
	for i := 0; i < len(ad.fields); i = i + 2 {
		pair := "\"" + ad.fields[i] + "\":\"" + ad.fields[i+1] + "\""
		if i == 0 {
			pairs = pair
			continue
		}
		pairs = pairs + ", " + pair
	}
	if len(pairs) > 0 {
		content = content + "\t{" + pairs + "}"
	}
	ad.mu.RUnlock()

	if _, err := fmt.Fprintf(ad.writer, "%04d-%02d-%02d %02d:%02d:%02d.%3d %-5s %s:%d\t%s\n",
		year, month, day, hour, min, sec, millsec, tag, file, line, content); err != nil {
		fmt.Printf("[ERROR] printf log error:" + err.Error())
		return err
	}
	return nil
}
