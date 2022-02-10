package main

import (
	"memtable/pkg/skiplist"
)

func main() {
	// exercise.Run(skiplist.New())

	// a0 := skiplist.NewNode(skiplist.NewKey([]byte{1}), nil)
	// a1 := skiplist.NewNode(skiplist.NewKey([]byte{1}), nil)
	// a1.Child = a0
	// b0 := skiplist.NewNode(skiplist.NewKey([]byte{2}), nil)
	// c0 := skiplist.NewNode(skiplist.NewKey([]byte{3}), nil)
	// d0 := skiplist.NewNode(skiplist.NewKey([]byte{4}), nil)
	// e0 := skiplist.NewNode(skiplist.NewKey([]byte{5}), nil)
	// e1 := skiplist.NewNode(skiplist.NewKey([]byte{5}), nil)
	// e1.Child = e0

	a0 := skiplist.NewNode(skiplist.NewKey([]byte{1}), nil)
	b0 := skiplist.NewNode(skiplist.NewKey([]byte{2}), nil)
	c0 := skiplist.NewNode(skiplist.NewKey([]byte{3}), nil)
	d0 := skiplist.NewNode(skiplist.NewKey([]byte{4}), nil)
	e0 := skiplist.NewNode(skiplist.NewKey([]byte{5}), nil)

	l := skiplist.New()
	l.Root().Append(a0).Append(b0).Append(c0).Append(d0).Append(e0)
	a1 := l.Lift(a0)
	e1 := l.Lift(e0)
	l.Lift(a1)
	l.Lift(e1)
	l.Print()

}
