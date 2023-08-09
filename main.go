package main

import (
	"fmt"
	"geecache"
	"net/http"
)

var db = map[string]string{
	"jie": "jie",
	"wen": "wen",
	"hao": "hao",
}

func main() {
	geecache.NewGroup("name", geecache.CacheCallBackFunc(
		func(key string) ([]byte, error) {
			if val, ok := db[key]; ok {
				return []byte(val), nil
			}
			return []byte{}, fmt.Errorf("not find")
		}), 2<<10)

	addr := "localhost:9999"
	peers := geecache.NewHttpPool(addr)
	http.ListenAndServe(addr, peers)
}
