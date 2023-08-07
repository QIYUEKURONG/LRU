package LRU

import "net/http"

func main() {

	addr := "localhost:9999"

	http.ListenAndServe(addr)
}
