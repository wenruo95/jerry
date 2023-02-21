/*================================================================
*   Copyright (C) 2021. All rights reserved.
*
*   file : password.go
*   coder: zemanzeng
*   date : 2021-10-06 10:57:01
*   desc : 密码生成
*
================================================================*/

package utils

import "math/rand"

func GeneratePassword(n int) string {
	if n <= 0 {
		return ""
	}

	lower := "abcdefghijklmnopqrstuvwxyz"
	upper := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers := "0123456789"
	special := "~!@#$%^&*-=_+,.?"

	s := lower + upper + numbers + special
	runes := make([]rune, 0, n)
	for i := 0; i < n; i++ {
		index := rand.Int() % len(s)
		runes = append(runes, rune(s[index]))
	}
	return string(runes)
}
