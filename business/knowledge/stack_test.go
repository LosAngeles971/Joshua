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
	s.push(p1)
	s.push(p2)
	s.push(p3)
	s.push(p4)
	if len(ss) != 4 || s.size() != 4 {
		t.FailNow()
	}
	s.pop()
	s.pop()
	s.pop()
	s.pop()
	if len(ss) != 4 || s.size() != 0 {
		t.FailNow()
	}
}