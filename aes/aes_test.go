/*================================================================
*   Copyright (C) 2021. All rights reserved.
*
*   file : aes_test.go
*   coder: zemanzeng
*   date : 2021-05-25 16:24:35
*   desc : aes测试用例
*
================================================================*/

package aes

import (
	"encoding/base64"
	"encoding/hex"
	"math/rand"
	"testing"
	"time"
)

func Test_AesAvailable(t *testing.T) {
	strs := []string{
		"hello world",
		"zemanzeng",
	}

	keys := []string{
		"abcdefghijklmnop",                      // 16位
		"abcdefghijklmnop" + "abcdefghijklmnop", // 32位
	}
	for _, key := range keys {
		for i := 0; i < len(strs); i++ {
			s := strs[i]
			bts, err := AesCBCEncrypt([]byte(s), []byte(key))
			if err != nil {
				t.Errorf("aes encrypt:%s key:%s error:%s", s, key, err.Error())
			}
			org, err := AesCBCDecrypt(bts, []byte(key))
			if err != nil {
				t.Errorf("aes desrypt:%s key:%v error:%s", string(bts), key, err.Error())
			}
			if string(org) != s {
				t.Errorf("aes result:%s dismatch:%s", string(org), s)
			}

			b16 := hex.EncodeToString(bts)
			b64 := base64.StdEncoding.EncodeToString(bts)
			t.Logf("s:%v key:%v aes_b16(%v):%v aes_b64(%v):%v", s, key, len(b16), b16, len(b64), b64)
		}
	}

}

func TestGenarateAesKey(t *testing.T) {

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 10; i++ {
		key := GenarateAesKey(AesKeyL16)
		t.Logf("index:%v len:%v key:%v", i, len(key), key)
	}
	for i := 0; i < 10; i++ {
		key := GenarateAesKey(AesKeyL32)
		t.Logf("index:%v len:%v key:%v", i, len(key), key)
	}
	for i := 0; i < 10; i++ {
		key := GenarateAesKey(AesKeyLRand)
		t.Logf("index:%v len:%v key:%v", i, len(key), key)
	}

}
