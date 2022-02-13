<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [Part 1.](#part-1)
- [Part 2.](#part-2)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Part 1.

Begin by designing an SSTable file format that you will use for storing key/value data in sorted order. There many details you may want to consider; as a starting point:

* Are you using a binary or text format?
* Are entries fixed size or variable size?
* How will you separate entries?
* What metadata will you store (e.g. length prefixes, sparse indexes)?
* What’s the process of finding a particular entry, and how efficient is it?

Once you’re satisfied with your design (you may also want to document it carefully as a long comment), add the following method to your in-memory key/value database:

```
// Flush the contents of the in-memory key/value database
// to `w` in the form of an SSTable.
flushSSTable(w *io.Writer) error
```
Now that you’ve added a method for writing an SSTable, the next step is to add the corresponding functionality for reading from it.

Implement the following interface (a subset of DB) for reading from a single SSTable file. You might find it helpful to use io.ReadSeeker for working with the open file.

```
type ImmutableDB interface {
    // Get gets the value for the given key. It returns an error if the
    // DB does not contain the key.
    Get(key []byte) (value []byte, err error)

    // Has returns true if the DB contains the given key.
    Has(key []byte) (ret bool, err error)

    // RangeScan returns an Iterator (see below) for scanning through all
    // key-value pairs in the given range, ordered by key ascending.
    RangeScan(start, limit []byte) (Iterator, error)
}
```

## Part 2.

Create a combined type that implements the original DB. interface. It should include a memtable together with a list of zero or more SSTables (i.e. types backed by a single SSTable file that implement ImmutableDB).

Whenever incoming writes cause the memtable to reach a fixed threshold (say, 2 MB), call flushSSTable, clear the memtable, and append the SSTable to the list.

Handle reads by attempting to read from the memtable as well as all SSTables. In what order should you attempt reads?

Feel free to ignore Delete and RangeScan for now, as they will be part of the next prework.

At this point, you should have a working key/value database that involves both an in-memory part and an on-disk part. However, note the following limitations:

You will lose the contents of the memtable if your program crashes (we will address this issue when we discuss logging)
The number of SSTables could quickly grow out of control (we will address when we discuss compaction)
Stretch Goals.

Compare your SSTable format to LevelDB’s. What decisions did you make differently? What are the tradeoffs between your two approaches?
Write a program that can parse LevelDB’s SSTable format: given an .ldb file, print out all keys and values stored in that file, in sorted order.
