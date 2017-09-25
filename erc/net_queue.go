package erc

import (
	"github.com/cirbo-lang/cirbo/cbo"
)

// netQueue is a FIFO queue of cbo.Net, implemented as a ring buffer
type netQueue struct {
	buf   []*cbo.Net
	start int
	end   int
	in    map[*cbo.Net]struct{}
}

func newNetQueue(capacity int) *netQueue {
	return &netQueue{
		buf:   make([]*cbo.Net, capacity+1),
		start: 0,
		end:   0,
		in:    make(map[*cbo.Net]struct{}, capacity),
	}
}

// Take removes the next item in the queue and returns it, or returns nil if
// the queue is empty.
func (q *netQueue) Take() *cbo.Net {
	ret := q.Peek()
	if ret != nil {
		q.start++
		delete(q.in, ret)
		if q.start == len(q.buf) {
			q.start = 0
		}
	}
	return ret
}

// Peek returns the next item in the queue without removing it. Returns nil
// if the queue is empty.
func (q *netQueue) Peek() *cbo.Net {
	if q.Len() == 0 {
		return nil
	}

	return q.buf[q.start]
}

func (q *netQueue) Has(n *cbo.Net) bool {
	_, in := q.in[n]
	return in
}

// Append adds the given net to the queue if it is not already present.
func (q *netQueue) Append(n *cbo.Net) {
	if q.Has(n) {
		return
	}

	if length := q.Len(); length == len(q.buf)-1 {
		// Need to grow our buffer
		newBuf := make([]*cbo.Net, length*2)
		length = copy(newBuf, q.buf[q.start:])
		length += copy(newBuf[length:], q.buf[:q.end])
		q.buf = newBuf
		q.start = 0
		q.end = length
	}

	q.buf[q.end] = n
	q.in[n] = struct{}{}
	q.end++
	if q.end == len(q.buf) {
		q.end = 0
	}
}

// Len returns the number of items in the queue
func (q *netQueue) Len() int {
	if q.end >= q.start {
		return q.end - q.start
	}

	return len(q.buf) - q.start - q.end
}
