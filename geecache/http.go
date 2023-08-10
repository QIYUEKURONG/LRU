package geecache

import (
	"fmt"
	"geecache/consistenthash"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

var DefaultPath = "/_geecache/"
var defaultReplicas = 20

type HttpPool struct {
	slft        string
	basePath    string
	peers       *consistenthash.HashMap
	mu          sync.Mutex
	httpPeesMap map[string]*HttpGetter
}

func NewHttpPool(self string) *HttpPool {
	return &HttpPool{slft: self, basePath: DefaultPath}
}

func (p *HttpPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.slft, fmt.Sprintf(format, v...))
}

func (p *HttpPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.peers = consistenthash.NewHashMap(nil, defaultReplicas)
	p.peers.Add(peers...)
	p.httpPeesMap = make(map[string]*HttpGetter, len(peers))

	for _, peer := range peers {
		p.httpPeesMap[peer] = &HttpGetter{
			baseURL: peer + p.basePath,
		}
	}
}

func (p *HttpPool) PeerPicker(key string) (PeerGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if perr := p.peers.Get(key); perr != "" && perr != p.slft {
		p.Log("pick peer %s", perr)
		return p.httpPeesMap[perr], true
	}
	return nil, false
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

type HttpGetter struct {
	baseURL string
}

func (p *HttpGetter) Log(format string, v ...interface{}) {
	log.Printf("[HttpGetter Server %s] %s", fmt.Sprintf(format, v...))
}

func (p *HttpGetter) Get(group string, key string) ([]byte, error) {
	requestUrl := fmt.Sprintf("%v%v/%v", p.baseURL, url.QueryEscape(group), url.QueryEscape(key))
	p.Log("url ---> %v", requestUrl)

	resultResponse, err := http.Get(requestUrl)
	if err != nil {
		p.Log("%v", err)
		return nil, err
	}
	defer resultResponse.Body.Close()

	if resultResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HttpGetter htttp request find error")
	}

	bytes, err := ioutil.ReadAll(resultResponse.Body)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

// can check if object that come true the interface all method
var _ PeerGetter = (*HttpGetter)(nil)
