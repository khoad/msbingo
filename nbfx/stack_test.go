package nbfx

import "testing"

func TestStack(t *testing.T) {
	stack := Stack{}

	assertEqual(t, stack.Len(), 0)
	stack.Push("hello")
	assertEqual(t, stack.Len(), 1)
	stack.Push("world")
	assertEqual(t, stack.Len(), 2)

	val := stack.Pop()
	assertEqual(t, val, "world")
	assertEqual(t, stack.Len(), 1)
	val = stack.Pop()
	assertEqual(t, val, "hello")
	assertEqual(t, stack.Len(), 0)

	val = stack.Pop()
	assertEqual(t, val, nil)
	assertEqual(t, stack.Len(), 0)
}
