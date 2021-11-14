package knowledge

import (
	"testing"
)

func TestStack(t *testing.T) {
	p1 := &Path{
		Path: []*Edge{},
	}
	p2 := &Path{
		Path: []*Edge{},
	}
	p3 := &Path{
		Path: []*Edge{},
	}
	p4 := &Path{
		Path: []*Edge{},
	}
	ss := []*Path{p1, p2, p3, p4}
	s := Stack{}
	s.Push(p1)
	s.Push(p2)
	s.Push(p3)
	s.Push(p4)
	if len(ss) != 4 || s.Size() != 4 {
		t.FailNow()
	}
	s.Pop()
	s.Pop()
	s.Pop()
	s.Pop()
	if len(ss) != 4 || s.Size() != 0 {
		t.FailNow()
	}
}