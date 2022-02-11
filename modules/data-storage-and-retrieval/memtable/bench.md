## Put
Benchmarks for inserting 1 element into various-sized skiplists:
```
goos: darwin
goarch: amd64
pkg: memtable/pkg/skiplist
cpu: Intel(R) Core(TM) i7-4870HQ CPU @ 2.50GHz
BenchmarkPut4k-8         2506168               407.7 ns/op
BenchmarkPut8k-8         3239389               437.2 ns/op
BenchmarkPut16k-8        3282692               386.1 ns/op

```

## Get
MAX_LEVEL=16, 1/5 chance of leveling:
```
goos: darwin
goarch: amd64
pkg: memtable/pkg/skiplist
BenchmarkGet4k-8         2283508               478.9 ns/op
BenchmarkGet8k-8         2155312               513.7 ns/op
BenchmarkGet16k-8        1856026               636.3 ns/op
BenchmarkGet32k-8        1752372               689.9 ns/op
```
