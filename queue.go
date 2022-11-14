package uq

const size = 64 // initial ring buffer capacity

/*
	queue
*/

type Queue[T any] struct {
	first, last *list[T]
}

func (q *Queue[T]) Push(val T) {
	// empty?
	if q.first == nil {
		next := new(list[T]).Push(val)
		q.first = next
		q.last = next
	} else {
		q.last = q.last.Push(val) // push to back
	}
}

func (q *Queue[T]) Shift() (val T, ok bool) {
	val, q.first, ok = q.first.Shift() // pop from front
	return
}

/*
	linked list
*/

// list is a singly-linked list.  Following the chain of 'next' pointers
// provides FIFO semantics.
type list[T any] struct {
	next  *list[T]
	queue rbuf[T]
}

func (l *list[T]) Push(val T) *list[T] {
	if l.queue.Push(val) {
		return l
	}

	l.next = l.grow().Push(val)
	return l.next // return the last element
}

func (l *list[T]) Shift() (t T, next *list[T], ok bool) {
	if l != nil {
		// l not empty OR l is head?
		if t, ok = l.queue.Shift(); ok {
			next = l
		} else if t, next, ok = l.next.Shift(); next == nil {
			next = l.shrink()
		}
	}

	return
}

func (l *list[T]) grow() *list[T] {
	return &list[T]{
		queue: rbuf[T]{
			// Increase array capacity geometrically, with r = 2.
			// This gives us amortized O(1) performance.
			array: make([]T, l.queue.Cap()<<1),
		},
	}
}

func (l *list[T]) shrink() *list[T] {
	// Only shrink if not at the minimum size.
	if l.queue.Cap() != size {
		l.queue = rbuf[T]{
			array: make([]T, l.queue.Cap()>>1),
		}
	}

	return l
}

/*
	ring buffer
*/

type rbuf[T any] struct {
	read, write uint32 // uint32 ensures both fit in single cache line
	array       []T
}

func (rb *rbuf[T]) Len() int {
	return int(rb.write - rb.read)
}

func (rb *rbuf[T]) Cap() int {
	return len(rb.array)
}

func (rb *rbuf[T]) Empty() bool {
	if rb.array == nil {
		rb.array = make([]T, size)
	}

	return rb.read == rb.write
}

func (rb *rbuf[T]) Full() bool {
	if rb.array == nil {
		rb.array = make([]T, size)
	}

	return rb.Len() == len(rb.array)
}

func (rb *rbuf[T]) Push(val T) (ok bool) {
	if ok = !rb.Full(); ok {
		rb.write++
		rb.array[rb.write&rb.mask()] = val
	}

	return
}

func (rb *rbuf[T]) Shift() (val T, ok bool) {
	if ok = !rb.Empty(); ok {
		rb.read++
		val = rb.array[rb.read&rb.mask()]
	}

	return
}

func (rb *rbuf[T]) mask() uint32 {
	return uint32(len(rb.array)) - 1
}
