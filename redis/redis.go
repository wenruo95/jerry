/*================================================================
*   Copyright (C) 2021. All rights reserved.
*
*   file : redis.go
*   coder: zemanzeng
*   date : 2021-04-12 15:33:18
*   desc : redis客户端初始化
*
================================================================*/

package redis

import (
	"errors"
	"fmt"
	"sync"

	"github.com/go-redis/redis"
)

var redisInsts = new(sync.Map)

// RedisSource redis连接相关配置
type RedisSource struct {
	Source     string `json:"source"`
	Passwd     string `json:"passwd"`
	SelectedDB int    `json:"selected_db"`
	PoolSize   int    `json:"pool_size"`
}

func GetRedis(source RedisSource) (*redis.Client, error) {
	if len(source.Source) == 0 {
		return nil, errors.New("invalid source")
	}
	if inst, exist := redisInsts.Load(source.Source); exist {
		return inst.(*redis.Client), nil
	}

	client := redis.NewClient(&redis.Options{
		Addr:     source.Source,
		Password: source.Passwd,
		DB:       source.SelectedDB,
		PoolSize: source.PoolSize,
	})
	redisInsts.Store(source.Source, client)

	if _, err := client.Ping().Result(); err != nil {
		return client, fmt.Errorf("ping error:%w", err)
	}
	return client, nil
}
