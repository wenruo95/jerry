/*================================================================
*   Copyright (C) 2021. All rights reserved.
*
*   file : error_code_test.go
*   coder: zemanzeng
*   date : 2021-05-26 18:46:07
*   desc : error_code测试用例
*
================================================================*/

package utils

import (
	"fmt"
	"testing"
)

func Test_ErrorCodeAvailable(t *testing.T) {
	RegisterErrorCode(100, "code_100")
	RegisterErrorCode(101, "code_101")
	RegisterErrorCode(102, "code_102")
	RegisterErrorCode(103, "code_103")
	RegisterErrorCode(104, "code_104")
	RegisterErrorCode(105, "code_105")

	arr := []int32{100, 101, 102, 103, 104, 105, 106}
	for _, code := range arr {
		//err := ECError(code, "", fmt.Sprintf("dididiext:%v", code))
		err := ECErrorExt(code, fmt.Sprintf("dididiext:%v", code))
		if code2 := GetErrorCode(err); code2 != code {
			t.Errorf("code:%v not equal:%v", code, code2)
		}
		t.Logf("code:%v code2:%v msg:%v ext:%v error:%v",
			code, GetErrorCode(err), GetErrorMsg(err), GetErrorExtMsg(err), err.Error())
	}

}
