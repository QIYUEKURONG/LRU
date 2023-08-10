package geecache

// httpPool can select different service by trans hash value of key

type PeerPicker interface {
	PeerPicker(key string) (PeerGetter, bool)
}

// client can get data by cache
// PeerGetter interface can get data by http process

type PeerGetter interface {
	Get(group string, key string) ([]byte, error)
}
