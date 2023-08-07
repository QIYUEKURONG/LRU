package cache

import (
	"net/http"
	"strings"
)

var DefaultPath = "/_geecache/"

type HttpPool struct {
	slft     string
	basePath string
}

func (h HttpPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if !strings.HasPrefix(r.URL.Path, h.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}
	// /<basepath>/<groupname>/<key> required
	parts := strings.SplitN(r.URL.Path[len(h.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

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
