/*
Instructions
============

Consider the following interface for an “ID service”:

	type idService interface {
    // Returns values in ascending order; it should be safe to call
    // getNext() concurrently without any additional synchronization.
    getNext() uint64
	}

Implement this interface using each of the following four strategies:

- Don’t perform any synchronization
- Atomically increment a counter value using sync/atomic
- Use a sync.Mutex to guard access to a shared counter value
- Launch a separate goroutine with exclusive access to a private counter value;
	handle getNext() calls by making “requests” and receiving “responses” on two
	separate channels.
*/
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"go.uber.org/atomic"
)

type idService interface{ getNext() uint64 }

type noSync struct{ n uint64 }

func (s *noSync) getNext() uint64 {
	ret := s.n
	s.n = s.n + 1
	return ret
}

type atomicInc struct{ n atomic.Uint64 }

func (a *atomicInc) getNext() uint64 {
	ret := a.n.Load()
	a.n.Inc()
	return ret
}

type lockedInc struct {
	n uint64
	m sync.Mutex
}

func (a *lockedInc) getNext() uint64 {
	a.m.Lock()
	defer a.m.Unlock()

	ret := a.n
	a.n = a.n + 1
	return ret
}

type channelInc struct {
	req chan (struct{})
	res chan (uint64)
}

func (c *channelInc) getNext() uint64 {
	c.req <- struct{}{}
	select {
	case n := <-c.res:
		return n
	case <-time.After(time.Second):
		panic("timeout!")
	}
}

func getSomeIDs(n int, getter idService) []uint64 {
	ids := make([]uint64, n)
	for i := 0; i < n; i++ {
		ids[i] = getter.getNext()
		time.Sleep(time.Duration(rand.Int63n(5)))
	}
	return ids
}

func run(n int, g idService) [][]uint64 {
	results := make([][]uint64, n)
	wg := sync.WaitGroup{}
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			results[i] = getSomeIDs(10, g)
			wg.Done()
		}(i)
	}
	wg.Wait()
	return results
}

func main() {
	n := 15
	// g := &noSync{}
	// g := &atomicInc{}
	// g := &lockedInc{}
	var count uint64 = 1
	req := make(chan (struct{}))
	res := make(chan (uint64))
	g := &channelInc{req: req, res: res}
	go func() {
		for {
			<-req
			count = count + 1
			res <- count
		}
	}()
	arr := run(n, g)
	fmt.Printf("%v\n", arr)
}
