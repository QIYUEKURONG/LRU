package cache

import (
	"fmt"
	"testing"
)

var db = map[string]string{
	"jie": "jie",
	"wen": "wen",
	"hao": "hao",
}

func TestAdd(t *testing.T) {
	group := NewGroup("name", CacheCallBackFunc(
		func(key string) ([]byte, error) {
			if val, ok := db[key]; ok {
				return []byte(val), nil
			}
			return []byte{}, fmt.Errorf("not find")
		}), 2<<10)

	//test1
	val, err := group.GetKey("jie")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(val)
	val, err = group.GetKey("jie")
	fmt.Println(val)
}
