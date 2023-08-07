package cache

import (
	"cache/lru"
	"sync"
)

//lru cache's Encapsulation layer
//this layer will provide add and get function to chat whit user

type cache struct {
	mutex         sync.Mutex
	lru           *lru.Cache
	cacheCapacity int64
}

func (c *cache) add(key string, v ByteView) {
	if key == "" {
		panic("key is nil")
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.lru == nil {
		c.lru = lru.NewLRU(c.cacheCapacity, nil)
	}
	c.lru.Add(key, v)
}

func (c *cache) get(key string) (ByteView, bool) {
	if key == "" {
		panic("key is nil")
	}
	if c.lru == nil {
		return ByteView{}, false
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	v, isHave := c.lru.Get(key)
	return v.(ByteView), isHave
}
