/*================================================================
*   Copyright (C) 2022. All rights reserved.
*
*   file : setup_test.go
*   coder: zemanzeng
*   date : 2022-01-20 15:04:32
*   desc : setup test case
*
================================================================*/

package log

import (
	"os"
	"testing"
)

func TestRegister(t *testing.T) {
	writer := NewWriter(os.Stdout)
	RegisterWithWriter("testcase1", "debug", writer)
	Get("testcase1").Debug("testcase1 ==> console")

	if err := RegisterWithFileName("testcase2", "debug", "test.log"); err != nil {
		t.Errorf("register_with_file_name error:" + err.Error())
		return
	}
	Get("testcase2").Debug("testcase2 ==> test.log")

	cl := NewAdapter("debug", GetWriter("testcase2"), 0)
	cl.SetPrefix("TESTCASE")
	cl.Debug("test prefix ==> test.log")

	cfg := NewLogConfig("debug", "jerry.log")
	if err := RegisterWithLogConfig(DefaultLogName, cfg); err != nil {
		t.Errorf("register_with_log_config error:" + err.Error())
	}

	Info("hello world ==> jerry")
}
