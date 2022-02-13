<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [Part 1.](#part-1)
- [Part 2.](#part-2)
- [Stretch Goals](#stretch-goals)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Part 1.

Begin by implementing the interface given below (which is a simplified version of LevelDB’s interface).

Note that RangeScan must return key-value pairs in sorted order.
For this first implementation, you don’t need to worry about efficiency; just get something working however you can (for example, it’s acceptable to have a linear time Get, or to sort on every Put call).
```go
type DB interface {
    // Get gets the value for the given key. It returns an error if the
    // DB does not contain the key.
    Get(key []byte) (value []byte, err error)

    // Has returns true if the DB contains the given key.
    Has(key []byte) (ret bool, err error)

    // Put sets the value for the given key. It overwrites any previous value
    // for that key; a DB is not a multi-map.
    Put(key, value []byte) error

    // Delete deletes the value for the given key.
    Delete(key []byte) error

    // RangeScan returns an Iterator (see below) for scanning through all
    // key-value pairs in the given range, ordered by key ascending.
    RangeScan(start, limit []byte) (Iterator, error)
}

type Iterator interface {
    // Next moves the iterator to the next key/value pair.
    // It returns false if the iterator is exhausted.
    Next() bool

    // Error returns any accumulated error. Exhausting all the key/value pairs
    // is not considered to be an error.
    Error() error

    // Key returns the key of the current key/value pair, or nil if done.
    Key() []byte

    // Value returns the value of the current key/value pair, or nil if done.
    Value() []byte
}
```

## Part 2.

Optimize your implementation. You may use any efficient (faster than linear time for Get and Put) in-memory data structure of your choice, but we’d strongly recommend using a Skip List instead of, say, a red-black tree because (1) it’s much easier to implement, (2) it’s much easier to add concurrency, and (3) it’s the data structure used in LevelDB. That being said, the rest of the project will work regardless of the data structure you choose, and you would definitely also learn a lot from choosing an alternative such as a red-black tree, AVL tree, B tree, or treap.

If you decide to use a Skip List, you may want to consult [Skip Lists: A Probabilistic Alternative to Balanced Trees](https://www.epaperpress.com/sortsearch/download/skiplist.pdf) as well as this [skip list visualization](https://people.ok.ubc.ca/ylucet/DS/SkipList.html).
Before you begin your optimized implementation, you may want to add tests (correctness will be much more difficult compared to Part 1) and benchmarks (to show off how much faster your optimized implementation is compared to Part 1).

## Stretch Goals

Compare your implementation to an existing implementation. If you used a Skip List, you might consider [LevelDB](https://github.com/google/leveldb/blob/master/db/skiplist.h), [goleveldb](https://github.com/syndtr/goleveldb/blob/master/leveldb/memdb/memdb.go), or [Redis](https://github.com/redis/redis/blob/1c71038540f8877adfd5eb2b6a6013a1a761bc6c/src/t_zset.c).
Consider how to support concurrent access from multiple goroutines. Can you do this without locking the entire structure for each access?
