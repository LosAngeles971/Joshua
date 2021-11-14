package state

import (
	"testing"
)

func TestAddElements(t *testing.T) {
	data := map[string]float64{
		"A": 1.0,
		"B": 2.0,
		"C": 3.0,
		"D": 4.0,
	}
	s := New()
	for k, v := range data {
		s.Add(k, v)
	}
	for k, v := range data {
		vv, ok := s.Get(k)
		if !ok {
			t.Fatalf("variable %v missing", k)
		}
		dd, ok := vv.GetValue()
		if !ok {
			t.Fatalf("variable %v is not defined", k)
		}
		if dd.(float64) != v {
			t.Fatalf("variable %v should have value %v not %v", k, v, dd)
		}
	}
}

func TestIsSubsetOf(t *testing.T) {
	c1 := State{
		Data: map[string]*Variable{},
	}
	c2 := State{
		Data: map[string]*Variable{},
	}
	c1.Add("a", 1.0)
	c1.Add("b", -1.0)
	c2.Add("a", 1.0)
	c2.Add("b", -1.0)
	c2.Add("c", 0.0)
	if ok := c1.IsSubsetOf(c2); !ok {
		t.Error("c1 must be part of c2!")
	}
	if ok := c2.IsSubsetOf(c1); ok {
		t.Error("c2 cannot be part of c1!")
	}
}

func TestClone(t *testing.T) {
	c1 := State{
		Data: map[string]*Variable{},
	}
	c1.Add("a", 1.0)
	c1.Add("b", -1.0)
	c1.Declare("c", TYPE_BOOL)
	c2 := c1.Clone()
	if ok := c1.IsSubsetOf(c2); !ok {
		t.Error("c1 must be part of c2!")
	}
	if ok := c2.IsSubsetOf(c1); !ok {
		t.Error("c2 must be part of c1!")
	}
}