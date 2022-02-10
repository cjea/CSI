## Put
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

## Get
MAX_LEVEL=16, 1/5 chance of leveling:
```
goos: darwin
goarch: amd64
pkg: memtable/pkg/skiplist
BenchmarkGet4k-8         2218915               541.6 ns/op
BenchmarkGet8k-8         1670294               656.1 ns/op
BenchmarkGet16k-8        1443441               823.2 ns/op
BenchmarkGet32k-8        1380772               866.6 ns/op
```
