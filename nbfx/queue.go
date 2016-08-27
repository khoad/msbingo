package nbfx

// Queue is a struct with "first", "last" and "length" properties
//https://gist.github.com/bemasher/1777766
type Queue struct {
	first  *queueElement
	last   *queueElement
	length int
}

type queueElement struct {
	value interface{} // All types satisfy the empty interface, so we can store anything here.
	next  *queueElement
}

// Len returns the queue's length
func (q *Queue) Len() int {
	return q.length
}

// Enqueue adds a new element into the queue
func (q *Queue) Enqueue(value interface{}) {
	element := &queueElement{value, nil}
	if q.first == nil {
		q.first = element
		q.last = element
	} else {
		q.last.next = element
		q.last = element
	}
	q.length++
}

// Dequeue removes the first element from the queue and return its value
// If the queue is empty, return nil
func (q *Queue) Dequeue() (value interface{}) {
	if q.length > 0 {
		value, q.first = q.first.value, q.first.next
		q.length--
		return
	}
	return nil
}
