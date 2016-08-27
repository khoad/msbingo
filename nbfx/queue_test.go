package nbfx

import "testing"

func TestQueue(t *testing.T) {
	queue := Queue{}

	assertEqual(t, queue.Len(), 0)
	queue.Enqueue("hello")
	assertEqual(t, queue.Len(), 1)
	queue.Enqueue("world")
	assertEqual(t, queue.Len(), 2)

	val := queue.Dequeue()
	assertEqual(t, val, "hello")
	assertEqual(t, queue.Len(), 1)
	val = queue.Dequeue()
	assertEqual(t, val, "world")
	assertEqual(t, queue.Len(), 0)

	val = queue.Dequeue()
	assertEqual(t, val, nil)
	assertEqual(t, queue.Len(), 0)
}
