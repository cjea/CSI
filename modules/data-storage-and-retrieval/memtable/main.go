package main

import (
	"memtable/pkg/exercise"
	"memtable/pkg/storage/skiplist"
)

func main() {
	exercise.Run(skiplist.New())
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
