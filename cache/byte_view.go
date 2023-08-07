package cache

type ByteView struct {
	bytes []byte
}

func (b ByteView) Len() int {
	return len(b.bytes)
}

func Clone(item []byte) []byte {
	result := make([]byte, len(item))
	copy(result, item)
	return result
}
