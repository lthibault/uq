package uq

const (
	size = 64       // ring buffer capacity
	mask = size - 1 // ring buffer slot-mask
)

/*
	queue
*/

type Queue[T any] struct {
	first, last *list[T]
}

func (q *Queue[T]) Push(val T) {
	// empty?
	if q.first == nil {
		next := new(list[T]).Push(val) // TODO:  pool.Get().Push(val)
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

	l.next = new(list[T]) // TODO:  pool.Get()
	_ = l.next.Push(val)  // always succeeds
	return l.next         // return the last element
}

func (l *list[T]) Shift() (T, *list[T], bool) {
	// l not empty OR l is head?
	if val, ok := l.queue.Shift(); ok || l.next == nil {
		return val, l, ok
	}

	// TODO: pool.Put(l)
	defer l.Reset()

	return l.next.Shift()
}

func (l *list[T]) Reset() {
	l.next = nil
	l.queue.Reset()
}

/*
	ring buffer
*/

type rbuf[T any] struct {
	read, write uint32 // uint32 ensures both fit in single cache line
	array       [size]T
}

func (rb *rbuf[T]) Len() int {
	return int(rb.write - rb.read)
}

func (rb *rbuf[T]) Empty() bool {
	return rb.read == rb.write
}

func (rb *rbuf[T]) Full() bool {
	return rb.Len() == size
}

func (rb *rbuf[T]) Push(val T) (ok bool) {
	if ok = !rb.Full(); ok {
		rb.write++
		rb.array[rb.write&mask] = val
	}

	return
}

func (rb *rbuf[T]) Shift() (val T, ok bool) {
	if ok = !rb.Empty(); ok {
		rb.read++
		val = rb.array[rb.read&mask]
	}

	return
}

func (rb *rbuf[T]) Reset() {
	rb.read = 0
	rb.write = 0
}
