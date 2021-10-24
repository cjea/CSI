Implement the following functions:
​
Given a string, return its length (number of bytes) without using len.
Given a type Point struct { x int, y int }, return its y coordinate without using p.y.
Given an []int, return the sum of values without using range or [].
Given a map[int]int, return the max value, again without using range or [].
​
The goal of this implementation exercise is for you to gain familiarity with the underlying representations of these basic types (hence the strange constraints).
​
You will find it helpful to use uintptr and unsafe.Pointer (see 13.1 - 13.2 of The Go Programming Language). You will likely also want to consult the Go source code, especially runtime/map.go, which you can find at /usr/local/go/src or on Github.
