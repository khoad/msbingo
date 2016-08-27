package nbfx

// Stack is a struct with "top" and "size" properties
// https://gist.github.com/bemasher/1777766
type Stack struct {
	top  *stackElement
	size int
}

// StackElement can be used in Stack struct
type stackElement struct {
	value interface{} // All types satisfy the empty interface, so we can store anything here.
	next  *stackElement
}

// Len returns the stack's length
func (s *Stack) Len() int {
	return s.size
}

// Push a new element onto the stack
func (s *Stack) Push(value interface{}) {
	s.top = &stackElement{value, s.top}
	s.size++
}

// Pop removes the top element from the stack and return its value
// If the stack is empty, return nil
func (s *Stack) Pop() (value interface{}) {
	if s.size > 0 {
		value, s.top = s.top.value, s.top.next
		s.size--
		return
	}
	return nil
}
