package consistenthash

import (
	"sort"
	"strconv"
)

type hash func([]byte) int

type HashMap struct {
	//product hash value by hashFunc
	hashFunc hash
	//prevent most request in single node
	replicas int
	// virtual list
	keys []int
	//key virtual node value real node
	hashMap map[int]string
}

func NewHashMap(hashFunc hash, replicas int) *HashMap {
	return &HashMap{
		hashFunc: hashFunc,
		replicas: replicas,
		hashMap:  make(map[int]string),
	}
}

func (h *HashMap) Add(keys ...string) {
	if len(keys) == 0 {
		return
	}
	for _, key := range keys {
		for i := 0; i < h.replicas; i++ {
			hashV := h.hashFunc([]byte(strconv.Itoa(i) + key))
			h.keys = append(h.keys, hashV)
			h.hashMap[hashV] = key
		}
	}
	sort.Ints(h.keys)
}

func (h *HashMap) Get(key string) string {
	if key == "" {
		return ""
	}
	hashV := h.hashFunc([]byte(key))

	idx := sort.Search(len(h.keys), func(i int) bool {
		return h.keys[i] >= hashV
	})

	return h.hashMap[h.keys[idx%len(h.keys)]]
}
