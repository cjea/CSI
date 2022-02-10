package exercise

import (
	"fmt"
	"memtable/pkg/storage"
)

func Run(db storage.DB) {
	for i := 0; i < 26; i++ {
		must(db.Put([]byte{byte(i + 97)}, []byte(fmt.Sprintf("val_%c", i+97))))
	}
	db.Print()
	it, err := db.RangeScan([]byte("ba"), []byte("z"))
	must(err)
	for it.Next() {
		fmt.Printf("Val: '%s'\n", string(it.Value()))
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
