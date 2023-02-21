/*================================================================
*   Copyright (C) 2021. All rights reserved.
*
*   file : business_auth_test.go
*   coder: zemanzeng
*   date : 2021-09-06 11:06:33
*   desc : 测试用例
*
================================================================*/

package token

import (
	"encoding/json"
	"testing"
)

type testBizMessage struct {
	Data string `json:"data"`
}

func (m *testBizMessage) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

func (m *testBizMessage) Unmarshal(buff []byte) error {
	return json.Unmarshal(buff, m)
}

var (
	S1 uint32 = 123
	S2 uint32 = 456
)

func (m *testBizMessage) GetSecretKey(secretId uint32) string {
	switch secretId {
	case S1:
		return "hello_world_0123"
	case S2:
		return "hello_world_0456"
	}
	return ""
}

func Test_BusinessAuthAvailable(t *testing.T) {

	m := &testBizMessage{Data: "hello world zemanzeng"}
	token, expire, err := Encode(S1, m, 24*60*60)
	if err != nil {
		t.Errorf("encode error:" + err.Error())
	}
	t.Logf("encode succ. token:%v length:%v expire:%v", token, len(token), expire)

	m2 := &testBizMessage{}
	if _, err := Decode(token, m2); err != nil {
		t.Errorf("decode token error:" + err.Error())
	}
	t.Logf("decode succ. info2:%+v", m2)

	if m.Data != m2.Data {
		t.Errorf("info not equals. info:%+v info2:%+v", m, m2)
	}
}

func Benchmark_BusinessAuthEncode(t *testing.B) {
	m := &testBizMessage{Data: "hello world zemanzeng"}

	for i := 0; i < t.N; i++ {
		_, _, err := Encode(S1, m, 24*60*60)
		if err != nil {
			t.Errorf("encode error:" + err.Error())
		}
	}

}

func Benchmark_BusinessAuthDecode(t *testing.B) {
	m := &testBizMessage{Data: "hello world zemanzeng"}
	token, _, err := Encode(S1, m, 24*60*60)
	if err != nil {
		t.Errorf("encode error:" + err.Error())
		return
	}

	for i := 0; i < t.N; i++ {
		m2 := &testBizMessage{}
		_, err := Decode(token, m2)
		if err != nil {
			t.Errorf("decode error:" + err.Error())
		}
	}

}
