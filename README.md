# uq
Fast unbounded queue with efficient allocation

## Benchmarks

Below are some benchmarks for the `Push` and `Shift` operations.  Benchmarks are notoriously difficult to interpret,
but these suggest that `uq.Queue` will not be a bottleneck in most applications.

The queue is implemented as a linked list of ring buffers, each with a capacity of 64 items, which amortizes heap
allocations.  Moreover, when the queue is emptied, the final ring buffer is not released, so that a subsequen call
to `Push` does not immediately allocate a new one.  These optimizations result in constant-average-time allocs in
applications with moderate loads.

```
go test -benchmem -run=^$ -bench ^BenchmarkQueue$ github.com/lthibault/uq

goos: darwin
goarch: amd64
pkg: github.com/lthibault/uq
cpu: Intel(R) Core(TM) i7-1068NG7 CPU @ 2.30GHz
BenchmarkQueue/Push-8         	148702400	         8.103 ns/op	       9 B/op	       0 allocs/op
BenchmarkQueue/Shift-8        	250869399	         4.375 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/lthibault/uq	7.999s
```
