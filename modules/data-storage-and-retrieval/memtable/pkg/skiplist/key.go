package skiplist

type Key struct {
	key []byte
	Raw []byte
}

func (k Key) Eq(k2 Key) bool {
	b1 := k.key
	b2 := k2.key
	if len(b1) != len(b2) {
		return false
	}
	for i := range b1 {
		if b1[i] != b2[i] {
			return false
		}
	}
	return true
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
