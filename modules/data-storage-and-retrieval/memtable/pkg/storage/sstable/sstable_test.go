package sstable

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	table := []struct {
		data    []byte
		numKeys uint32
		err     error
	}{
		{
			data:    []byte{0xDB, 0xDB, 0xDB, 1, 0, 0, 0},
			numKeys: 1,
			err:     nil,
		},
		{
			data:    []byte{0xDB, 0xDB, 0xDB, 0, 1, 0, 0},
			numKeys: 256,
			err:     nil,
		},
		{
			data:    []byte{0xDB, 0xDB, 0xDB, 4, 3, 2, 1},
			numKeys: 16909060,
			err:     nil,
		},
		{
			data:    []byte{0xDB, 0xDB},
			numKeys: 1,
			err:     ErrMissingDbMagicNum,
		},
		{
			data:    []byte{0xDB, 0xFF, 0xDB},
			numKeys: 1,
			err:     ErrMissingDbMagicNum,
		},
		{
			data:    []byte{0xDB, 0xDB, 0xDB, 0, 1, 0},
			numKeys: 1,
			err:     ErrMissingNumKeys,
		},
	}

	for i, tt := range table {
		name := fmt.Sprintf("test-%d", i)
		t.Run(name, func(t *testing.T) {
			res, err := NewParser(tt.data)
			if tt.err != err {
				t.Errorf("expected error %v; got %v", tt.err, err)
			}
			if err != nil {
				return
			}
			if res.NumKeyIdxs != tt.numKeys {
				t.Errorf("expected numKeys %v; got %v", tt.numKeys, res.NumKeyIdxs)
			}
		})
	}
}
