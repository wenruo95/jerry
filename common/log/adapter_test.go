/*================================================================
*   Copyright (C) 2022. All rights reserved.
*
*   file : adapter_test.go
*   coder: zemanzeng
*   date : 2022-01-20 12:38:31
*   desc : adapter_test
*
================================================================*/

package log

import (
	"os"
	"testing"
)

func TestAdapterAvailable(t *testing.T) {
	ad := NewAdapter("debug", NewWriter(os.Stdout), 0)
	ad.SetPrefix("TestCase")
	ad.Trace("debug level trace log")
	ad.Tracef("level:%v tracef log", ad.GetLevel())

	ad.Debug("debug level debug log")
	ad.Debugf("level:%v debugf log", ad.GetLevel())

	ad.Info("debug level info log")
	ad.Infof("level:%v infof log", ad.GetLevel())

	ad.Warn("debug level warn log")
	ad.Warnf("level:%v warnf log", ad.GetLevel())

	ad.Error("debug level error log")
	ad.Errorf("level:%v errorf log", ad.GetLevel())

	ad2 := ad.WithFields("hello", "world", "test1")
	ad2.Trace("debug level trace log")
	ad2.Tracef("level:%v tracef log", ad2.GetLevel())

	ad2.Debug("debug level debug log")
	ad2.Debugf("level:%v debug log", ad2.GetLevel())

	ad2.Info("debug level info log")
	ad2.Infof("level:%v info log", ad2.GetLevel())

	ad2.Warn("debug level warn log")
	ad2.Warnf("level:%v warn log", ad2.GetLevel())

	ad2.Error("debug level error log")
	ad2.Errorf("level:%v error log", ad2.GetLevel())
}

func TestOutput(t *testing.T) {
	ad := NewAdapter("debug", NewWriter(os.Stdout), 0)
	ad.Output(0, "INFO", "hello world deep=0")
	ad.Output(1, "INFO", "hello world deep=1")
	ad.Output(2, "INFO", "hello world deep=2")
	ad.Output(3, "INFO", "hello world deep=3")

	// 2022-01-20 14:50:22.843 INFO  log/adapter.go:189	hello world deep=0
	// 2022-01-20 14:50:22.843 INFO  log/adapter_test.go:56	hello world deep=1
	// 2022-01-20 14:50:22.843 INFO  testing/testing.go:1259	hello world deep=2
	// 2022-01-20 14:50:22.843 INFO  runtime/asm_amd64.s:1581	hello world deep=3
}
