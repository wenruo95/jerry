/*================================================================
*   Copyright (C) 2021. All rights reserved.
*
*   file : pack_test.go
*   coder: zemanzeng
*   date : 2021-02-04 05:03:38
*   desc : pack测试用例
*
================================================================*/

package tcp

import (
	"bytes"
	"testing"
)

func Test_pack(t *testing.T) {
	data := make([]byte, 1024)
	buf := Pack(data, 1, 100)
	if len(buf) != len(data)+12 {
		t.Fatalf("invalid pack size. size:%v\n", len(buf))
		return
	}

	data1, messageFlag, txid, err := Unpack(bytes.NewBuffer(buf))
	if err != nil {
		t.Fatal(err)
		return
	}

	if messageFlag != 1 {
		t.Fatalf("message flag is not 1\n")
		return
	}

	if txid != 100 {
		t.Fatalf("txid is not 100\n")
		return
	}

	if len(data1) != len(data) {
		t.Fatalf("message size error\n")
		return
	}

	if bytes.Compare(data, data1) != 0 {
		t.Fatalf("data != data1")
		return
	}
}
