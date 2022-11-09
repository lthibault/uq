package uq

import (
	"testing"

	"github.com/stretchr/testify/require"
)

/*
	queue
*/

func TestQueue(t *testing.T) {
	t.Parallel()

	var q Queue[int]
	for i := 0; i < size*2; i++ {
		q.Push(i)
	}

	for i := 0; i < size*2; i++ {
		val, ok := q.Shift()
		require.Equal(t, i, val, "should have FIFO semantics")
		require.True(t, ok, "should report success")
	}

	val, ok := q.Shift()
	require.Zero(t, val, "should return zero value")
	require.False(t, ok, "should report empty")
}

func BenchmarkQueue(b *testing.B) {
	b.Run("Push", func(b *testing.B) {
		b.ReportAllocs()

		var q Queue[int]
		for i := 0; i < b.N; i++ {
			q.Push(i)
		}
	})

	b.Run("Shift", func(b *testing.B) {
		var q Queue[int]
		for i := 0; i < b.N; i++ {
			q.Push(i)
		}

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, _ = q.Shift()
		}
	})
}

/*
	linked list
*/

func TestList(t *testing.T) {
	t.Parallel()
	t.Helper()

	var (
		first = new(list[int])
		last  = first
	)

	t.Run("Push", func(t *testing.T) {
		for i := 0; i < size*2; i++ {
			last = last.Push(i)
			require.NotNil(t, last, "should not return nil list")
		}

		// Check invariants for last element
		require.Nil(t, last.next, "should return last element")
		require.True(t, last.queue.Full(), "last element should be full")

		// Check invariants for first element
		require.Equal(t, last, first.next, "first element should point to last")
		require.True(t, first.queue.Full(), "first element should be full")
	})

	t.Run("Shift", func(t *testing.T) {
		var val int
		var ok bool
		for i := 0; i < size*2; i++ {
			val, first, ok = first.Shift()
			require.Equal(t, i, val, "should have FIFO semantics")
			require.NotNil(t, first, "should not return nil list")
			require.True(t, ok, "should report success")
		}

		require.Nil(t, first.next, "should be head")
		require.True(t, first.queue.Empty(), "queue should be empty")
		require.Equal(t, last, first,
			"last element should have been promoted to first")
	})
}

func BenchmarkList(b *testing.B) {
	b.Run("Push", func(b *testing.B) {
		b.ReportAllocs()

		bufs := b.N / size
		if b.N%size != 0 {
			bufs++
		}

		b.ResetTimer()
		b.ReportMetric(float64(bufs), "links")

		var (
			first = new(list[int])
			last  = first
		)

		for i := 0; i < b.N; i++ {
			last = last.Push(i)
		}
	})

	b.Run("Shift", func(b *testing.B) {
		b.ReportAllocs()

		bufs := b.N / size
		if b.N%size != 0 {
			bufs++
		}

		var (
			first = new(list[int])
			last  = first
		)

		for i := 0; i < b.N; i++ {
			last = last.Push(i)
		}

		b.ResetTimer()
		b.ReportMetric(float64(bufs), "unlinks")

		for i := 0; i < b.N; i++ {
			_, first, _ = first.Shift()
		}
	})
}

/*
	ring buffer
*/

func TestRingBuffer(t *testing.T) {
	t.Parallel()
	t.Helper()

	var rb rbuf[int]

	t.Run("Push", func(t *testing.T) {
		for i := 0; i < size; i++ {
			require.True(t, rb.Push(i), "should not be full")
		}

		require.Equal(t, size, rb.Len(), "should have len=cap")
		require.True(t, rb.Full(), "should be full")
		require.False(t, rb.Push(9001),
			"should fail to push to full buffer")
	})

	t.Run("Shift", func(t *testing.T) {
		for i := 0; i < size; i++ {
			val, ok := rb.Shift()
			require.True(t, ok, "should not be empty")
			require.Equal(t, i, val, "should have FIFO semantics")
		}

		require.Zero(t, rb.Len(), "should have len=0")
		require.True(t, rb.Empty(), "should be empty")

		_, ok := rb.Shift()
		require.False(t, ok, "should fail to shift empty buffer")
	})
}

func BenchmarkRingBuffer(b *testing.B) {
	b.ReportAllocs()

	var rb rbuf[int]
	for i := 0; i < b.N; i++ {
		for rb.Push(i) {
		}

		for _, ok := rb.Shift(); ok; _, ok = rb.Shift() {
		}
	}

	b.StopTimer()
	b.ReportMetric(float64(size), "push/op")
	b.ReportMetric(float64(size), "shift/op")
}
