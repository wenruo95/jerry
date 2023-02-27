package utils

import (
	"sync"
	"time"

	"github.com/golang/groupcache/lru"
)

// LRUCache lru缓存with ttl
type LRUCache struct {
	cache *lru.Cache
	lock  sync.RWMutex
}

// Value interface封装 支持过期机制
type Value struct {
	value  interface{}
	expire int64
}

func NewLRUCache(bucket int) *LRUCache {
	return &LRUCache{
		cache: lru.New(bucket),
	}
}

func (cache *LRUCache) Set(key interface{}, value interface{}) {
	cache.cache.Add(key, &Value{value: value, expire: -1})
}

func (cache *LRUCache) SetWithTTL(key interface{}, value interface{}, ttl int64) {
	cache.cache.Add(key, &Value{value: value, expire: time.Now().Unix() + ttl})
}

func (cache *LRUCache) Get(key interface{}) (interface{}, bool) {
	if value, exist := cache.cache.Get(key); exist {
		if v2, ok := value.(*Value); ok {
			if v2.expire < 0 || v2.expire > time.Now().Unix() {
				return v2.value, true
			}
		}
	}
	return nil, false
}

func (cache *LRUCache) GetString(key interface{}) string {
	if v, exist := cache.Get(key); exist {
		if v2, ok := v.(string); ok {
			return v2
		}
	}
	return ""
}
