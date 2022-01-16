/*================================================================
*   Copyright (C) 2022. All rights reserved.
*
*   file : async_writer.go
*   coder: zemanzeng
*   date : 2022-01-16 18:14:53
*   desc :
*
================================================================*/

package log

import (
	"bytes"
	"errors"
	"io"
	"time"
)

// AsyncWriter 日志异步写入类 实现zap.WriteSyncer接口
type AsyncWriter struct {
	cfg *AsyncConfig

	writer io.Writer

	logQueue chan []byte
	syncChan chan struct{}
}

// NewAsyncWriter 根据传入的参数创建一个RollWriter对
func NewAsyncWriter(writer io.Writer, cfg *AsyncConfig) *AsyncWriter {
	w := &AsyncWriter{
		writer:   writer,
		cfg:      cfg,
		logQueue: make(chan []byte, cfg.QueueSize),
		syncChan: make(chan struct{}),
	}

	// 批量异步写入日志
	go w.batchWriteLog()
	return w
}

// Write 实现io.Writer
func (w *AsyncWriter) Write(data []byte) (int, error) {
	log := make([]byte, len(data))
	copy(log, data)
	if w.cfg.DropLog {
		select {
		case w.logQueue <- log:
		default:
			return 0, errors.New("log queue is full")
		}
	} else {
		w.logQueue <- log
	}
	return len(data), nil
}

// Sync 实现zap.WriteSyncer接口
func (w *AsyncWriter) Sync() error {
	w.syncChan <- struct{}{}
	return nil
}

// Close 实现io.Closer
func (w *AsyncWriter) Close() error {
	return w.Sync()
}

// batchWriteLog 批量异步写入日志
func (w *AsyncWriter) batchWriteLog() {
	dur := time.Millisecond * time.Duration(w.cfg.IntervalMs)
	buffer := bytes.NewBuffer(make([]byte, 0, w.cfg.WriteSize*2))
	timer := time.NewTimer(dur)
	for {
		select {
		case <-timer.C:
			if buffer.Len() > 0 {
				_, _ = w.writer.Write(buffer.Bytes())
				buffer.Reset()
			}
			timer.Reset(dur)
		case data := <-w.logQueue:
			buffer.Write(data)
			if buffer.Len() >= w.cfg.WriteSize {
				_, _ = w.writer.Write(buffer.Bytes())
				buffer.Reset()
			}
		case <-w.syncChan:
			if buffer.Len() > 0 {
				_, _ = w.writer.Write(buffer.Bytes())
				buffer.Reset()
			}
			size := len(w.logQueue)
			for i := 0; i < size; i++ {
				v := <-w.logQueue
				_, _ = w.writer.Write(v)
			}
		}
	}
}
