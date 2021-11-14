package knowledge

/*
Stack implements a stack data structure of Path.
It is useful to check a set of possible paths,
processing the elements with a LIFO approach.
*/
type Stack []*Path

func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

func (s *Stack) Size() int {
	return len(*s)
}

func (s *Stack) Push(p *Path) {
	*s = append(*s, p)
}

func (s *Stack) Pop() (*Path, bool) {
	if s.IsEmpty() {
		return nil, false
	} else {
		index := len(*s) - 1
		element := (*s)[index]
		*s = (*s)[:index]
		return element, true
	}
}