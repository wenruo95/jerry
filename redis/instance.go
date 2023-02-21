/*================================================================
*   Copyright (C) 2021. All rights reserved.
*
*   file : instance.go
*   coder: zemanzeng
*   date : 2021-09-29 16:31:23
*   desc : redis_instance
*
================================================================*/

package redis

import "github.com/go-redis/redis"

type Instance struct {
	*redis.Client
	prefix string
}

func NewInstance(cli *redis.Client, prefix string) *Instance {
	return &Instance{
		Client: cli,
		prefix: prefix,
	}
}

func (ri *Instance) Key(biz string, subKey string) string {
	return ri.prefix + ":" + biz + ":" + subKey
}
