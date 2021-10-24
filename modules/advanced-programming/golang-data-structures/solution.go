package main

import (
	"fmt"
	"unsafe"
)

const INT_SIZE = unsafe.Sizeof(int(1))
const POINTER_SIZE = unsafe.Sizeof(&struct{}{})

type Point struct {
	x int
	y int
}

// Given a type Point struct { x int, y int }, return its y coordinate without using p.y.
func (p Point) yCoord() int {
	pStart := unsafe.Pointer(&p)
	yPtr := uintptr(pStart) + INT_SIZE
	return *(*int)(unsafe.Pointer(yPtr))
}

// Given a string, return its length (number of bytes) without using len.
func myLength(s string) int {
	start := unsafe.Pointer(&s)
	lengthPtr := (uintptr(start) + POINTER_SIZE)
	return *(*int)(unsafe.Pointer(lengthPtr))
}

// Given an []int, return the sum of values without using range or [].
func mySum(ns []int) int {
	sum := 0
	slicePtr := uintptr(unsafe.Pointer(&ns))
	firstDatumPtr := *(**int)(unsafe.Pointer(slicePtr))

	l := *(*int)(unsafe.Pointer(slicePtr + POINTER_SIZE))
	datumPtr := uintptr(unsafe.Pointer(firstDatumPtr))
	for i := 0; i < l; i++ {
		n := *(*int)(unsafe.Pointer(datumPtr))
		sum += n
		datumPtr += INT_SIZE
	}
	return sum
}

// Given a map[int]int, return the max value, again without using range or [].
func maxVal(m map[int]int) int {
	// TODO
	return 0
}

func main() {
	str := "abc"
	p := Point{x: 1, y: 5}
	nums := []int{1, 2, 3}
	m := map[int]int{1: 3, 2: 2, 3: 1}

	fmt.Printf(
		"myLength(%#v) = %d\n%#v.yCoord() = %d\nmySum(%#v) = %d\nmaxVal(%#v) = %d",
		str, myLength(str), p, p.yCoord(), nums, mySum(nums), m, maxVal(m),
	)
}
