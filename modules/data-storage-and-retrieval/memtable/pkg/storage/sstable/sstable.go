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
	- Num keys: 8 bytes
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

type SSTable struct {
	Data    []byte
	NumKeys uint64
}

func validateTable(data []byte) error {
	start := *(*int)(unsafe.Pointer(&data[0]))
	if start != MAGIC_NUM_TABLE {
		return fmt.Errorf("corrupted data: expected MAGIC_NUM_TABLE")
	}
	return nil
}

func New(data []byte) (*SSTable, error) {
	if err := validateTable(data); err != nil {
		return nil, err
	}
	numKeys := *(*uint64)(unsafe.Pointer(&data[3]))
	return &SSTable{
		Data:    data,
		NumKeys: numKeys,
	}, nil
}

func Foo() {
	x := []byte{0xDB, 0xDB, 0xDB}
	err := validateTable(x)
	fmt.Printf("Err: %v\n", err)
}
