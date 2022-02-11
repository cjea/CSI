package skiplist

import "bytes"

type Key struct {
	key []byte
	Raw []byte
}

func (k Key) Eq(k2 Key) bool {
	cmp := bytes.Compare(k.key, k2.key)
	return cmp == 0
}

func (k Key) Lt(k2 Key) bool {
	cmp := bytes.Compare(k.key, k2.key)
	return cmp == -1
}

func (k Key) Lte(k2 Key) bool {
	cmp := bytes.Compare(k.key, k2.key)
	return cmp == 0 || cmp == -1
}

// Prepend a sentinel value to the input key, to ensure that all linked lists
// can start with a guaranteed minimum value of {0}.
func NewKey(b []byte) Key {
	ret := make([]byte, len(b)+1)
	ret[0] = 1
	for i := 0; i < len(b); i++ {
		ret[i+1] = b[i]
	}
	return Key{key: ret, Raw: b}
}
