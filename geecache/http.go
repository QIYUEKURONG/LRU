package geecache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

var DefaultPath = "/_geecache/"

type HttpPool struct {
	slft     string
	basePath string
}

func NewHttpPool(self string) *HttpPool {
	return &HttpPool{slft: self, basePath: DefaultPath}
}

func (p *HttpPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.slft, fmt.Sprintf(format, v...))
}

func (p HttpPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}
	p.Log("%s %s", r.Method, r.URL.Path)
	// /<basepath>/<groupname>/<key> required
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	vv := r.URL.Query()
	fmt.Println(vv)

	groupName := parts[0]
	key := parts[1]

	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}

	view, err := group.GetKey(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(view.bytes)
}
