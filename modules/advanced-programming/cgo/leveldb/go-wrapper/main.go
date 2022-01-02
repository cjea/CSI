package main

/*
#cgo CFLAGS: -I/usr/local/Cellar/leveldb/1.23/include
#cgo LDFLAGS: -L/usr/local/Cellar/leveldb/1.23/lib -l leveldb

#include "leveldb/c.h"
#include <stdlib.h>
*/
import "C"

import "fmt"

func main() {
	opts := C.leveldb_options_create()
	fmt.Printf("%v\n", opts)
}
