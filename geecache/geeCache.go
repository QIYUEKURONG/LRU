package geecache

import (
	"fmt"
	"sync"
)

type Group struct {
	Name      string
	CallBack  CacheCallBack
	MainCache *cache
}

type CacheCallBack interface {
	Get(string) ([]byte, error)
}

type CacheCallBackFunc func(string) ([]byte, error)

func (c CacheCallBackFunc) Get(key string) ([]byte, error) {
	return c(key)
}

var (
	mu    sync.RWMutex
	group map[string]*Group
)

func GetGroup(name string) *Group {
	mu.Lock()
	defer mu.Unlock()
	return group[name]
}

func NewGroup(name string, callBack CacheCallBack, MaxCacheCap int64) *Group {
	if callBack == nil {
		//if not have data, can use callBack to laod data
		panic("not have callBack")
	}
	mu.Lock()
	defer mu.Unlock()
	if item, ok := group[name]; ok {
		return item
	}

	return &Group{Name: name, CallBack: callBack, MainCache: &cache{cacheCapacity: MaxCacheCap}}
}

func (g *Group) GetKey(key string) (ByteView, error) {

	if key == "" {
		return ByteView{}, fmt.Errorf("parameter error")
	}

	if v, ok := g.MainCache.get(key); ok {
		return v, fmt.Errorf("cache not find key")
	}

	return g.loadCallBack(key)
}

func (g *Group) loadCallBack(key string) (ByteView, error) {
	if g.CallBack == nil {
		return ByteView{}, fmt.Errorf("not find Callbaack")
	}

	item, err := g.CallBack.Get(key)
	if err != nil {
		return ByteView{}, fmt.Errorf("not find about key")
	}

	result := ByteView{
		bytes: Clone(item),
	}

	g.loadCallBackToCache(key, result)

	return result, err

}

func (g *Group) loadCallBackToCache(key string, val ByteView) {
	g.MainCache.add(key, val)
}
