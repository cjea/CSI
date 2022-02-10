package exercise

import (
	"fmt"
	"memtable/pkg/storage"
	"memtable/pkg/storage/skiplist"
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

func Seed(n int) *skiplist.Skiplist {
	l := skiplist.New()
	for i := 0; i < n; i++ {
		str := fmt.Sprint(i)
		l.Put([]byte("key-"+str), []byte("val-"+str))
	}
	return l
}

var SeededList4k = Seed(1 << 12)
var SeededList8k = Seed(1 << 13)
var SeededList16k = Seed(1 << 14)

func must(err error) {
	if err != nil {
		panic(err)
	}
}
