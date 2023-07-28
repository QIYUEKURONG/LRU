package lru

import (
	"fmt"
	"testing"
)

type String string

func (s String) Len() int {
	return len(s)
}

func TestAdd(t *testing.T) {
	cache := NewLRU(2<<31, nil)
	cache.Add("score", String("123"))
	cache.Add("score1", String("456"))
	v, isHave := cache.Get("score1")
	if isHave {
		fmt.Println(v)
	}
}
