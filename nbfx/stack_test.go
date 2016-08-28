package nbfx

import "testing"

func TestStack(t *testing.T) {
	stack := stack{}

	assertEqual(t, stack.size, 0)
	stack.push("hello")
	assertEqual(t, stack.size, 1)
	stack.push("world")
	assertEqual(t, stack.size, 2)

	val := stack.pop()
	assertEqual(t, val, "world")
	assertEqual(t, stack.size, 1)
	val = stack.pop()
	assertEqual(t, val, "hello")
	assertEqual(t, stack.size, 0)

	val = stack.pop()
	assertEqual(t, val, nil)
	assertEqual(t, stack.size, 0)
}
