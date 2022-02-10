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
Retrieving 1 constant element from various-sized skiplists  (MAX_LEVEL=16, 1/2 chance of leveling):
```
goos: darwin
goarch: amd64
pkg: memtable/pkg/skiplist
BenchmarkGet4k-8         2772039               422.8 ns/op
BenchmarkGet8k-8         2938028               400.6 ns/op
BenchmarkGet16k-8        2494647               481.0 ns/op
```

Retrieving 1 random element from various-sized skiplists (MAX_LEVEL=16, 1/2 chance of leveling):
```
goos: darwin
goarch: amd64
pkg: memtable/pkg/skiplist
BenchmarkGet4k-8          646668              1667 ns/op
BenchmarkGet8k-8          692256              1612 ns/op
BenchmarkGet16k-8         507258              2512 ns/op
```
Retrieving 1 random element from various-sized skiplists (MAX_LEVEL=16, 1/4 chance of leveling):
```
goos: darwin
goarch: amd64
pkg: memtable/pkg/skiplist
BenchmarkGet4k-8          866635              1502 ns/op
BenchmarkGet8k-8          705526              1495 ns/op
BenchmarkGet16k-8         591199              2378 ns/op
```
