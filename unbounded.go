package unbounded

import (
	"sync"
)

// chunk represents a fixed-size chunk of items.
type chunk[T any] struct {
	items [128]T
	next  *chunk[T]
}

// waiter represents a node in the wait list.
type waiter[T any] struct {
	sync.Mutex
	dptr     *T
	next     *waiter[T]
	received bool
}

// unboundedChannel represents the unbounded channel with chunked buffers and a wait list.
type unboundedChannel[T any] struct {
	sync.Mutex
	head       *chunk[T]
	tail       *chunk[T]
	waitList   *waiter[T]
	lastWaiter *waiter[T]
	headIndex  int
	tailIndex  int
	closed     bool
}

// New creates a new unbounded channel.
func New[T any]() *unboundedChannel[T] {
	ch := &unboundedChannel[T]{
		head: &chunk[T]{},
		tail: &chunk[T]{},
	}
	ch.head = ch.tail
	return ch
}

// Send adds an item to the channel.
func (ch *unboundedChannel[T]) Send(value T) {
	ch.Lock()
	if ch.closed {
		ch.Unlock()
		panic("writing into closed channel")
	}

	if ch.waitList != nil {
		w := ch.waitList
		next := w.next
		if next != nil {
			ch.waitList = next
		} else {
			ch.waitList = nil
			ch.lastWaiter = nil
		}
		ch.Unlock()
		*w.dptr = value
		w.received = true
		w.Unlock()
		return
	}

	if ch.tailIndex == len(ch.tail.items) {
		newChunk := &chunk[T]{}
		ch.tail.next = newChunk
		ch.tail = newChunk
		ch.tailIndex = 0
	}

	ch.tail.items[ch.tailIndex] = value
	ch.tailIndex++
	ch.Unlock()
}

// Receive removes and returns an item from the channel.
func (ch *unboundedChannel[T]) Receive() (T, bool) {
	ch.Lock()
	if ch.headIndex == ch.tailIndex && ch.head == ch.tail && !ch.closed {
		// reading in wait list
		var value T // keep here for performance reason
		w := &waiter[T]{dptr: &value}
		w.Lock()
		if ch.lastWaiter == nil {
			ch.waitList = w
		} else {
			ch.lastWaiter.next = w
		}
		ch.lastWaiter = w
		ch.Unlock()
		w.Lock()
		return value, w.received
	} else {
		// reading from buffer
		var value T         // keep here for performance reason
		ind := ch.headIndex // register suggestion
		if ch.closed && ind == ch.tailIndex && ch.head == ch.tail {
			ch.Unlock()
			return value, false // zero value
		}
		value = ch.head.items[ind]
		ind++
		if ind == len(ch.head.items) {
			if ch.head.next != nil {
				ch.head = ch.head.next
			}
			ind = 0
		}
		ch.headIndex = ind
		ch.Unlock()
		return value, true
	}
}

// Close closes the channel.
func (ch *unboundedChannel[T]) Close() {
	ch.Lock()
	defer ch.Unlock()

	ch.closed = true

	for ch.waitList != nil {
		waiter := ch.waitList
		ch.waitList = ch.waitList.next
		waiter.Unlock()
	}
}
