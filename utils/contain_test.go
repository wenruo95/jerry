/*================================================================
*   Copyright (C) 2021. All rights reserved.
*
*   file : contain_test.go
*   coder: zemanzeng
*   date : 2021-09-27 18:35:13
*   desc : contain 测试用例
*
================================================================*/

package utils

import "testing"

func Test_StrListEqual(t *testing.T) {
	a := []string{"a", "b"}
	b := []string{"a", "b"}
	c := []string{"a", "b", "c"}
	if !StrListEqual(a, b, true) {
		t.Errorf("str list not equal? a:%+v b:%+v", a, b)
	}
	if StrListEqual(a, c, true) {
		t.Errorf("str list equal? a:%+v c:%+v", a, c)
	}
}

func Test_StrListMD5(t *testing.T) {
	a := []string{"a", "b"}
	b := []string{"b", "a"}
	c := []string{"a", "b", "c"}
	if StrListMD5(a, true) != StrListMD5(b, true) {
		t.Errorf("str list not equal? a:%+v b:%+v", a, b)
	}
	if StrListMD5(a, true) == StrListMD5(c, true) {
		t.Errorf("str list equal? a:%+v c:%+v", a, c)
	}
}
