package nbfx

// https://gist.github.com/bemasher/1777766
type queue struct {
	first  *queueElement
	last   *queueElement
	length int
}

type queueElement struct {
	value interface{} // All types satisfy the empty interface, so we can store anything here.
	next  *queueElement
}

func (q *queue) enqueue(value interface{}) {
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

func (q *queue) dequeue() (value interface{}) {
	if q.length > 0 {
		value, q.first = q.first.value, q.first.next
		q.length--
		return
	}
	return nil
}
