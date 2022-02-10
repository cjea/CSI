Benchmarks for inserting 1 element into various-sized skiplists:
```
goos: darwin
goarch: amd64
pkg: memtable/pkg/skiplist
cpu: Intel(R) Core(TM) i7-4870HQ CPU @ 2.50GHz
BenchmarkPut4k-8         2823948               425.5 ns/op
BenchmarkPut8k-8         2938836               495.7 ns/op
BenchmarkPut16k-8        2402826               554.4 ns/op

```
Retrieving 1 element from various-sized skiplists:
```
goos: darwin
goarch: amd64
pkg: memtable/pkg/skiplist
BenchmarkGet4k-8         2772039               422.8 ns/op
BenchmarkGet8k-8         2938028               400.6 ns/op
BenchmarkGet16k-8        2494647               481.0 ns/op
```
