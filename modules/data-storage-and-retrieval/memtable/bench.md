Benchmarks for inserting 1 element into various-sized skiplists:
```
goos: darwin
goarch: amd64
pkg: memtable/pkg/skiplist
cpu: Intel(R) Core(TM) i7-4870HQ CPU @ 2.50GHz
BenchmarkPut4k-8         2598684               460.1 ns/op
BenchmarkPut8k-8         2829456               482.9 ns/op
BenchmarkPut16k-8        2450562               507.5 ns/op

```
Retrieving 1 element from various-sized skiplists:
```
goos: darwin
goarch: amd64
pkg: memtable/pkg/skiplist
BenchmarkGet4k-8         1758028               679.7 ns/op
BenchmarkGet8k-8         1727550               704.3 ns/op
BenchmarkGet16k-8        1489754               798.7 ns/op
```
