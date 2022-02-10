Benchmarks for inserting 1 element into various-sized skiplists:
```
goos: darwin
goarch: amd64
pkg: memtable/pkg/skiplist
cpu: Intel(R) Core(TM) i7-4870HQ CPU @ 2.50GHz
BenchmarkPut4k-8         2493964               428.5 ns/op
BenchmarkPut8k-8         2631334               419.5 ns/op
BenchmarkPut16k-8        2487915               472.5 ns/op

```
Retrieving 1 element from various-sized skiplists:
```
goos: darwin
goarch: amd64
pkg: memtable/pkg/skiplist
BenchmarkGet4k-8         1988679               611.5 ns/op
BenchmarkGet8k-8         1780476               665.3 ns/op
BenchmarkGet16k-8        1523884               785.9 ns/op
```
