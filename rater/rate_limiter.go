/*================================================================
*   Copyright (C) 2020. All rights reserved.
*
*   file : rate_limiter.go
*   coder: zemanzeng
*   date : 2020-06-23 10:05:20
*   desc : 限频接口封装
*
================================================================*/

package rater

import (
	"golang.org/x/time/rate"
)

var defaltRateLimiter *RateLimiter

func InitRateLimiter(limit int) {
	defaltRateLimiter = NewRateLimit(limit)
}

func Allow() bool {
	return defaltRateLimiter.Allow()
}

type RateLimiter struct {
	rateLimiter *rate.Limiter
}

func NewRateLimit(limit int) *RateLimiter {
	rateLimiter := rate.NewLimiter(rate.Limit(limit), limit)

	return &RateLimiter{
		rateLimiter: rateLimiter,
	}
}

func (limiter *RateLimiter) Allow() bool {
	if limiter == nil || limiter.rateLimiter == nil {
		return true
	}

	return limiter.rateLimiter.Allow()
}
