package nbfx

import "testing"

func TestQueue(t *testing.T) {
	queue := queue{}

	assertEqual(t, queue.length, 0)
	queue.enqueue("hello")
	assertEqual(t, queue.length, 1)
	queue.enqueue("world")
	assertEqual(t, queue.length, 2)

	val := queue.dequeue()
	assertEqual(t, val, "hello")
	assertEqual(t, queue.length, 1)
	val = queue.dequeue()
	assertEqual(t, val, "world")
	assertEqual(t, queue.length, 0)

	val = queue.dequeue()
	assertEqual(t, val, nil)
	assertEqual(t, queue.length, 0)
}
