package nbfx

// https://gist.github.com/bemasher/1777766
type stack struct {
	top  *stackElement
	size int
}

type stackElement struct {
	value interface{} // All types satisfy the empty interface, so we can store anything here.
	next  *stackElement
}

func (s *stack) push(value interface{}) {
	s.top = &stackElement{value, s.top}
	s.size++
}

func (s *stack) pop() (value interface{}) {
	if s.size > 0 {
		value, s.top = s.top.value, s.top.next
		s.size--
		return
	}
	return nil
}
