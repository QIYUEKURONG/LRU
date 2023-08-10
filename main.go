package main

import (
	"flag"
	"fmt"
	"geecache"
	"log"
	"net/http"
)

var db = map[string]string{
	"jie": "jie",
	"wen": "wen",
	"hao": "hao",
}

func SetServerCache(selfAdds string, adds []string, geeCache *geecache.Group) {
	httpPools := geecache.NewHttpPool(selfAdds)
	httpPools.Set(adds...)
	geeCache.Peers = httpPools
	log.Fatal(http.ListenAndServe(selfAdds[7:], httpPools))
}

func ServerSelectRequest(key string, geeCache *geecache.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			geeCache.GetKey()

		}))

}

func main() {
	geecache := geecache.NewGroup("name", geecache.CacheCallBackFunc(
		func(key string) ([]byte, error) {
			if val, ok := db[key]; ok {
				return []byte(val), nil
			}
			return []byte{}, fmt.Errorf("not find")
		}), 2<<10)

	port := 8000
	api := false
	flag.IntVar(&port, "port", 6666, "server")
	flag.BoolVar(&api, "api", false, "can start")

	portAddsMap := map[int]string{
		6666: "localhost:6666",
		7777: "localhost:7777",
		8888: "localhost:8888",
	}
	var adds []string
	for _, value := range portAddsMap {
		adds = append(adds, value)
	}

	if api {

	}

	SetServerCache(portAddsMap[port], adds, geecache)
}
