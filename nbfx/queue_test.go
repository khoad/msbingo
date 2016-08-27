package nbfx

import "testing"

func TestQueue(t *testing.T) {
	queue := Queue{}

	assertEqual(t, queue.length, 0)
	queue.Enqueue("hello")
	assertEqual(t, queue.length, 1)
	queue.Enqueue("world")
	assertEqual(t, queue.length, 2)

	val := queue.Dequeue()
	assertEqual(t, val, "hello")
	assertEqual(t, queue.length, 1)
	val = queue.Dequeue()
	assertEqual(t, val, "world")
	assertEqual(t, queue.length, 0)

	val = queue.Dequeue()
	assertEqual(t, val, nil)
	assertEqual(t, queue.length, 0)
}

func assertEqual(t *testing.T, actual, expected interface{}) {
	if expected != actual {
		t.Errorf("%v not equal to expected %v", actual, expected)
	}
}
