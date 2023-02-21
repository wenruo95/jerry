/*================================================================
*   Copyright (C) 2021. All rights reserved.
*
*   file : contain.go
*   coder: zemanzeng
*   date : 2021-09-27 18:21:01
*   desc : 数组等容器
*
================================================================*/

package utils

import (
	"sort"
	"strings"
)

func ListContain(list []string, item string) bool {
	for _, key := range list {
		if key == item {
			return true
		}
	}
	return false
}

// 数据新增/删除元素
func CmpStrList(src []string, dst []string) (addn []string, deln []string) {
	if len(src) == 0 {
		return dst, nil
	}
	if len(dst) == 0 {
		return nil, src
	}

	amap := make(map[string]int)
	for index, v := range src {
		if _, ok := amap[v]; ok {
			continue
		}
		amap[v] = index
	}
	bmap := make(map[string]int)
	for index, v := range dst {
		if _, ok := bmap[v]; ok {
			continue
		}
		bmap[v] = index
	}

	for index, item := range dst {
		if _, ok := amap[item]; ok {
			continue
		}
		if v, ok := bmap[item]; !ok || v != index {
			continue
		}
		addn = append(addn, item)
	}
	for index, item := range src {
		if _, ok := bmap[item]; ok {
			continue
		}
		if v, ok := amap[item]; !ok || v != index {
			continue
		}
		deln = append(deln, item)
	}

	return addn, deln
}

// 数组交集
func MixedList(src []string, dst []string) []string {
	amap := make(map[string]struct{})
	for _, item := range src {
		amap[item] = struct{}{}
	}

	list := make([]string, 0)
	bmap := make(map[string]int)
	if len(src) < len(dst) {
		for _, item := range src {
			if _, ok := bmap[item]; !ok {
				continue
			}
			if _, ok := amap[item]; ok {
				continue
			}
			list = append(list, item)
		}
		return list
	}

	for _, item := range dst {
		if _, ok := amap[item]; !ok {
			continue
		}
		if _, ok := bmap[item]; ok {
			continue
		}
		list = append(list, item)
	}
	return list
}

// 删除重复元素 reserveFirst保留第第一次出现的元素 否则保留最后一次出现的元素
func RemoveRepeatStrItem(items []string, reserveFirst bool) []string {
	m := make(map[string]int)
	for index, item := range items {
		if reserveFirst {
			if _, ok := m[item]; ok {
				continue
			}
		}
		m[item] = index
	}

	list := make([]string, 0, len(m))
	for index, item := range items {
		if v, ok := m[item]; ok && index == v {
			list = append(list, item)
		}
	}
	return list
}

// 字符串数组数据是否相等 忽略位置、重复元素
func StrListIsEqual(src []string, dst []string) bool {
	c1, c2 := CmpStrList(src, dst)
	return (len(c1) == 0) && (len(c2) == 0)
}

// 忽略位置
func StrListEqual(a []string, b []string, doSort bool) bool {
	if len(a) != len(b) {
		return false
	}
	if doSort {
		sort.Strings(a)
		sort.Strings(b)
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func StrListMD5(list []string, doSort bool) string {
	if doSort {
		sort.Strings(list)
	}
	return MD5(strings.Join(list, " "))
}
