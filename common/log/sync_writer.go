/*================================================================
*   Copyright (C) 2022. All rights reserved.
*
*   file : sync_writer.go
*   coder: zemanzeng
*   date : 2022-01-16 17:45:12
*   desc : 同步写writer
*
================================================================*/

package log

import (
	"io"

	"github.com/natefinch/lumberjack"
)

type SyncWriter struct {
	writer *lumberjack.Logger
}

func NewSyncWriter(writer *WriterConfig) *SyncWriter {
	return &SyncWriter{
		writer: &lumberjack.Logger{
			Filename:   writer.FileName,
			MaxSize:    writer.MaxSize,
			MaxBackups: writer.MaxBackups,
			MaxAge:     writer.MaxAge,
			Compress:   writer.Compress,
			LocalTime:  writer.LocalTime,
		},
	}
}

func (w *SyncWriter) Write(data []byte) (int, error) {
	return w.writer.Write(data)
}

func (w *SyncWriter) Close() error {
	return w.writer.Close()
}

func (w *SyncWriter) Sync() error {
	return nil
}

type WriterAdapter struct {
	writer io.Writer
}

func NewWriter(writer io.Writer) *WriterAdapter {
	return &WriterAdapter{writer: writer}
}

func (w *WriterAdapter) Write(data []byte) (int, error) {
	return w.writer.Write(data)
}

func (w *WriterAdapter) Close() error {
	return nil
}

func (w *WriterAdapter) Sync() error {
	return nil
}
