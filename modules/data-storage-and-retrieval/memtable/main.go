package main

import (
	"memtable/pkg/storage/sstable"
)

func main() {
	// exercise.Run(skiplist.New())
	sstable.Foo()
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
