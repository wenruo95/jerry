/*================================================================
*   Copyright (C) 2021. All rights reserved.
*
*   file : dirty_words_test.go
*   coder: zemanzeng
*   date : 2021-06-22 21:01:07
*   desc : 脏词测试用例
*
================================================================*/

package dirtywords

import (
	"strings"
	"testing"
)

var dirtyList = []string{
	"cao尼玛",
	"fuck",
	"tmd",
	"妈的戈壁",
	"傻屌",
	"傻逼",
	"妈的",
	"尼玛",
	"操",
	"操nm",
}

func Test_DirtyWords(t *testing.T) {
	dw := NewDirtyWords()
	//dw.ReBuild(dirtyList, t.Logf)
	dw.ReBuild(dirtyList, nil)

	md := StrListMD5(dirtyList, true) // 这里会把dirtyList排序，放到后面执行
	expectList := []string{
		"cao尼玛",
		"fuck",
		"tmd",
		"傻屌",
		"傻逼",
		"妈的",
		"尼玛",
		"操",
	}
	list := dw.List()
	if equal := StrListEqual(expectList, list, true); !equal {
		t.Errorf("tree not match expect. expect[%d]:%v result[%d]:%v",
			len(expectList), expectList, len(list), list)
	}

	s := "黑夜总会散去，黎明终要来临"
	if words := dw.Match(s, true); len(words) != 0 {
		t.Errorf("match error match:%v len:%v", words, len(words))
	}
	if result, _ := dw.Filter(s, nil); result != s {
		t.Errorf("filter error expect[%d]:%v result[%d]:%v", len(s), s, len(result), result)
	}

	newList := make([]string, 0)
	newList = append(newList, dirtyList...)
	newList = append(newList, "夜总会")
	if StrListMD5(newList, true) != md {
		dw.ReBuild(newList, nil)
	}
	if words := dw.Match(s, true); len(words) != 1 {
		t.Errorf("match error match:%v len:%v", words, len(words))
	}
	s1 := "黑*散去，黎明终要来临"
	if result, _ := dw.Filter(s, nil); result != s1 {
		t.Errorf("filter error expect[%d]:%v result[%d]:%v", len(s1), s1, len(result), result)
	}

	filterFunc := func(words string, match string) string {
		return strings.Repeat("#", len([]rune(match)))
	}
	s2 := "黑###散去，黎明终要来临"
	if result, _ := dw.Filter(s, filterFunc); result != s2 {
		t.Errorf("filter error expect[%d]:%v result[%d]:%v", len(s2), s2, len(result), result)
	}

}
func Test_DirtyWordsWithIgnores(t *testing.T) {
	dw := NewDirtyWordsWithIgnores([]rune{'\t', '\r', '\n', ' '})
	dirtyList = append(dirtyList, "夜总会")
	dw.ReBuild(dirtyList, nil)

	s := "黑夜\t\n总\r 会散去，黎明终\t\r\n要 来临"
	expect := "黑*散去，黎明终\t\r\n要 来临"
	if result, _ := dw.Filter(s, nil); result != expect {
		t.Errorf("filter error expect[%d]:%v result[%d]:%v", len(expect), expect, len(result), result)
	}
}
