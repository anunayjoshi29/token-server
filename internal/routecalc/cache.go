package routecalc

import (
	"sync"
	"time"
)

type CacheItem struct {
	value      []RouteResult
	expiration int64
}

type Cache struct {
	data map[string]CacheItem
	mu   sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string]CacheItem),
	}
}

func (c *Cache) Set(key string, value []RouteResult, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = CacheItem{
		value:      value,
		expiration: time.Now().Add(duration).Unix(),
	}
}

func (c *Cache) Get(key string) ([]RouteResult, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, found := c.data[key]
	if !found || time.Now().Unix() > item.expiration {
		return nil, false
	}
	return item.value, true
}
