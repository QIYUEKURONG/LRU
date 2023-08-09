package lru

import "container/list"

type Cache struct {
	MaxCacheCap   int64
	AlreadyUseCap int64
	List          *list.List
	CacheMap      map[string]*list.Element
	OnEvicted     func(key string, value Value)
}

type Value interface {
	Len() int
}

type entry struct {
	key   string
	value Value
}

func NewLRU(MaxCacheCap int64, OnEvicted func(key string, value Value)) *Cache {
	return &Cache{
		MaxCacheCap: MaxCacheCap,
		List:        list.New(),
		CacheMap:    make(map[string]*list.Element),
		OnEvicted:   OnEvicted,
	}
}

func (c *Cache) Add(key string, value Value) {
	if element, ok := c.CacheMap[key]; ok {
		//存在 就调整位置
		c.List.MoveToFront(element)
		kv := element.Value.(entry)
		c.AlreadyUseCap += int64(value.Len()) - int64(kv.value.Len())
		element.Value = kv
	} else {
		e := c.List.PushFront(entry{key: key, value: value})
		c.CacheMap[key] = e
		c.AlreadyUseCap += int64(len(key)) + int64(value.Len())
	}
	if c.AlreadyUseCap > c.MaxCacheCap {
		c.RemoveFromTailElemet()
	}
}

func (c *Cache) Get(key string) (value Value, isHave bool) {
	if element, ok := c.CacheMap[key]; ok {
		c.List.MoveToFront(element)
		kv := element.Value.(entry)
		return kv.value, true
	}
	return
}

func (c *Cache) RemoveFromTailElemet() {
	ele := c.List.Back()
	if ele != nil {
		c.List.Remove(ele)
		entry := ele.Value.(entry)
		delete(c.CacheMap, entry.key)
		c.AlreadyUseCap -= int64(len(entry.key)) + int64(entry.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(entry.key, entry.value)
		}
	}
}
