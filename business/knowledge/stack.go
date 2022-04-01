/*
This file provides functions to handle a Stack,
that is a LIFO queue of Path.
*/
package knowledge


type Stack []*Path

func (s *Stack) isEmpty() bool {
	return len(*s) == 0
}

func (s *Stack) size() int {
	return len(*s)
}

func (s *Stack) push(p *Path) {
	*s = append(*s, p)
}

func (s *Stack) pop() (*Path, bool) {
	if s.isEmpty() {
		return nil, false
	} else {
		index := len(*s) - 1
		element := (*s)[index]
		*s = (*s)[:index]
		return element, true
	}
}