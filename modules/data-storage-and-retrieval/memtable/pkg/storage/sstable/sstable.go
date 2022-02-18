package sstable

import (
	"fmt"
	"unsafe"
)

/*
	* SStables will be binary files consisting of key-value pairs.
	* Data is immutable.
	* Entries are dynamically sized.
	* Entry format needs to support:
		- Finding an entry by key
		- Finding a range of entries from a given key

	Format of one key-value entry:
	==============================
	 	- Magic number: 3 bytes
		- Key sizes ~2^10 bytes
		- Value sizes ~2^14 bytes
		- Total *metadata* size (bits): 24 + 24 + 10 = 54 bits =~ 8 bytes
		- Total max size: 17KB + metadata
	 ___________________________________________________
	|              |              |            |       |
	| MAGIC_NUMBER | TOTAL_LENGTH | KEY_LENGTH | ENTRY |
	|______________|______________|____________|_______|

	Format of key index:
	===================
		- Key length: 10 bits
		- Key 2^10 bytes (max)
		- Offset from end: 4 bytes
		- Total size (bits): 1KB + 4 bytes + 10 bits
	 ____________________________________
	|            |     |                 |
	| KEY_LENGTH | KEY | OFFSET_FROM_END |
	|____________|_____|_________________|

	Format of database:
	===================
	- Magic number: 3 bytes
	- Num keys: 4 bytes
	 __________________________________________________________
	|                   |          |               |          |
	| MAGIC_FILE_NUMBER | NUM_KEYS | KEY_INDEX ... | ENTRy... |
	|___________________|__________|__________________________|

	Finding an entry:
	=================
	for 1..NUM_KEYS
		read KEY (using KEY_LENGTH)
			if KEY matches target: return getEntry(offset)
			else repeat for next key
	Time complexity: 0(n) w.r.t the number and size of all keys
*/

const (
	MAGIC_NUM_TABLE = 0xDBDBDB
	MAGIC_NUM_ENTRY = 0xEAC4AB
)

var (
	ErrMissingNumKeys    = fmt.Errorf("corrputed: expected numKeys")
	ErrMissingDbMagicNum = fmt.Errorf("corrupted: expected MAGIC_NUM_TABLE")
)

type SSTableParser struct {
	Pos     uint64
	Data    []byte
	NumKeys uint32
}

func (s *SSTableParser) Read(p []byte) (n int, err error) {
	outIdx := 0
	offset := int(s.Pos)
	remaining := min(len(s.Data)-offset, len(p))
	for ; outIdx < remaining; outIdx++ {
		p[outIdx] = s.Data[outIdx+offset]
	}
	s.Pos = uint64(outIdx + offset)
	return outIdx, nil
}

func (s *SSTableParser) ParseMagicNumTable() error {
	out := make([]byte, 3)
	n, err := s.Read(out)
	if n != 3 || err != nil {
		return ErrMissingDbMagicNum
	}

	if !everyEl(out, 0xDB) {
		return ErrMissingDbMagicNum
	}
	return nil
}

func (s *SSTableParser) ParseNumKeys() error {
	out := make([]byte, 4)
	n, err := s.Read(out)
	if n != 4 || err != nil {
		return ErrMissingNumKeys
	}
	s.NumKeys = *(*uint32)(unsafe.Pointer(&out[0]))
	return nil
}

func NewParser(data []byte) (*SSTableParser, error) {
	ret := &SSTableParser{Data: data, Pos: 0}
	if err := ret.ParseMagicNumTable(); err != nil {
		return nil, err
	}
	if err := ret.ParseNumKeys(); err != nil {
		return nil, err
	}
	return ret, nil
}

func Foo() {
	// x := []byte{0xDB, 0xDB, 0xDB}
	// err := validateTable(x)
	// fmt.Printf("Err: %v\n", err)
	panic("Foo not implemented :)")
}

func everyEl(bs []byte, b byte) bool {
	for _, el := range bs {
		if el != b {
			return false
		}
	}
	return true
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
