# uq

[![Godoc Reference](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/lthibault/uq)
[![Go Report Card](https://goreportcard.com/badge/github.com/lthibault/uq?style=flat-square)](https://goreportcard.com/report/github.com/lthibault/uq)
[![tests](https://github.com/lthibault/uq/workflows/Go/badge.svg)](https://github.com/lthibault/uq/actions/workflows/go.yml)


Fast unbounded queue with O(1) amortized allocations.

>**WARNING:**  use caution with unbounded memory.  Use appropriate flow-control measures to prevent OOM failures.

## Benchmarks

Below are some benchmarks for the `Push` and `Shift` operations.  Benchmarks are notoriously difficult to interpret,
but these suggest that `uq.Queue` will not be a bottleneck in most applications.

The queue is implemented as a linked list of ring buffers, of geometrically-increasing capacity (r=2).  Moreover, when the queue is emptied, the final ring buffer is not released, so that a subsequen call to `Push` does not immediately allocate a new one.  These optimizations result in amortized constant-time allocations.  Effectively, `uq` converges on the appropriate size for a single ring-buffer for any constant workload.

A geometric growth factor of `r=2` is required to ensure ring buffer capacities are always powers of 2.  On average, this means about 40% of the largest ring buffer is empty.  To mitigate this waste, the ring buffer size is correspondingly shrunk when the queue is empty.  Note that pathological workloads exist, which may either prevent an under-utilized buffer from being downsized or cause thrashing.  Both issues can be mitigated by implementing appropriate flow-control mechanisms to smooth out dataflow.  Future work should focus on refactoring the ring buffers to support arbitrary capacities, at which point we can set the growth rate to a more conservative `r=1.5` and the shrink rate to an asymmetrical value of `r=.25`.

```
$ go test -benchmem -run='^$' -bench '^BenchmarkQueue$' github.com/lthibault/uq
goos: darwin
goarch: amd64
pkg: github.com/lthibault/uq
cpu: Intel(R) Core(TM) i7-1068NG7 CPU @ 2.30GHz
BenchmarkQueue/Push-8           197861551                7.076 ns/op          10 B/op          0 allocs/op
BenchmarkQueue/Shift-8          262964143                4.545 ns/op           0 B/op          0 allocs/op
PASS
ok      github.com/lthibault/uq 6.805s
```
