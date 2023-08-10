package geecache

type PeerPicker interface {
	PeerPicker(key string) (PeerGetter, bool)
}

type PeerGetter interface {
	Get(group string, key string) ([]byte, error)
}
