package sstable

import (
	"bytes"
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
  	- Magic num: 3 bytes
    - Key length: 10 bits
    - Key 2^10 bytes (max)
    - Offset from end: 4 bytes
    - Total max size: 3 bytes + 10 bits + 1KB + 4 bytes
   ______________________________________________________________
  |                         |            |     |                 |
  | MAGIC_KEY_IDX_NUMBER    | KEY_LENGTH | KEY | OFFSET_FROM_END |
  |_________________________|____________|_____|_________________|

  Format of database:
  ===================
  - Magic number: 3 bytes
  - Num key indices: 4 bytes
   __________________________________________________________
  |                   |          |               |          |
  | MAGIC_FILE_NUMBER | NUM_IDXS | KEY_INDEX ... | ENTRY... |
  |___________________|__________|__________________________|

  Finding an entry:
  =================
  for 1..NUM_KEYS
    read KEY (using KEY_LENGTH)
      if KEY matches target: return getEntry(offset)
      else repeat for next key
  Time complexity: 0(n) w.r.t the number and size of all keys
*/

var (
	MAGIC_NUM_TABLE   = []byte{0xDB, 0xDB, 0xDB}
	MAGIC_NUM_KEY_IDX = []byte{0x4E, 0x4E, 0x4E}
	MAGIC_NUM_ENTRY   = []byte{0xEA, 0xC4, 0xAB}
)

var (
	ErrMissingNumKeys        = fmt.Errorf("corrputed: expected numKeys")
	ErrMissingDbMagicNum     = fmt.Errorf("corrupted: expected MAGIC_NUM_TABLE")
	ErrMissingKeyIdxMagicNum = fmt.Errorf("corrupted: expected MAGIC_NUM_KEY_IDX")
)

type KeyIdx struct {
	KeyLenBytes int
	Key         []byte
	Offset      int
}

type SSTableParser struct {
	Pos        uint64
	Data       []byte
	NumKeyIdxs uint32
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

func (s *SSTableParser) ParseMagicNum(magic []byte) bool {
	length := len(magic)
	out := make([]byte, length)
	s.Read(out)

	return bytes.Compare(out, magic) == 0
}

func (s *SSTableParser) ParseMagicNumTable() error {
	if s.ParseMagicNum(MAGIC_NUM_TABLE) {
		return nil
	}
	return ErrMissingDbMagicNum
}

func (s *SSTableParser) ParseNumKeyIdxs() error {
	out := make([]byte, 4)
	n, err := s.Read(out)
	if n != 4 || err != nil {
		return ErrMissingNumKeys
	}
	s.NumKeyIdxs = *(*uint32)(unsafe.Pointer(&out[0]))
	return nil
}

func (s *SSTableParser) ParseKeyIdx() (*KeyIdx, error) {
	if !s.ParseMagicNum(MAGIC_NUM_KEY_IDX) {
		return nil, ErrMissingKeyIdxMagicNum
	}
	// TODO: parse key length, key, and offset into *KeyIdx{}
	return nil, nil
}

func NewParser(data []byte) (*SSTableParser, error) {
	ret := &SSTableParser{Data: data, Pos: 0}

	if err := ret.ParseMagicNumTable(); err != nil {
		return nil, err
	}

	if err := ret.ParseNumKeyIdxs(); err != nil {
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
