/*================================================================
*   Copyright (C) 2020. All rights reserved.
*
*   file : rate_limiter_test.go
*   coder: zemanzeng
*   date : 2020-11-12 03:26:04
*   desc : 限频测试用例
*
================================================================*/

package rater

import (
	"log"
	"testing"
	"time"
)

func TestRateLimiter(t *testing.T) {

	x := 10

	InitRateLimiter(x + 1)

	for i := 0; i < 100; i++ {
		log.Printf("i:%v", Allow())
		if i%(x*2) == 0 {
			time.Sleep(time.Second)
		}
	}

}
