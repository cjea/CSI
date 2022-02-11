## Put
Benchmarks for inserting 1 element into various-sized skiplists:
```
goos: darwin
goarch: amd64
pkg: memtable/pkg/skiplist
cpu: Intel(R) Core(TM) i7-4870HQ CPU @ 2.50GHz
BenchmarkPut1k-8         2515980               466.7 ns/op
BenchmarkPut2k-8         2093922               502.0 ns/op
BenchmarkPut4k-8         2449112               583.0 ns/op
BenchmarkPut8k-8         2025996               568.1 ns/op
BenchmarkPut16k-8        2064651               578.1 ns/op

```

## Get
MAX_LEVEL=16, 1/5 chance of leveling:
```
goos: darwin
goarch: amd64
pkg: memtable/pkg/skiplist
BenchmarkGet4k-8         1856839               640.7 ns/op
BenchmarkGet8k-8         1762884               673.6 ns/op
BenchmarkGet16k-8        1519165               785.2 ns/op
BenchmarkGet32k-8        1448210               823.1 ns/op
```
