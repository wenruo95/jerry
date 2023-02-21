/*================================================================
*   Copyright (C) 2021. All rights reserved.
*
*   file : password_test.go
*   coder: zemanzeng
*   date : 2021-10-06 11:04:50
*   desc : 密码测试用例
*
================================================================*/

package utils

import (
	"math/rand"
	"testing"
	"time"
)

func Test_GeneratePassword(t *testing.T) {

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 10; i++ {
		t.Logf("i:%v password:%v", i, GeneratePassword(15))
	}
}
