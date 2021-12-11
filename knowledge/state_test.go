package knowledge

import (
	"testing"
)

func TestDatamapCreation(t *testing.T) {
	data := map[string]float64{
		"A": 1.0,
		"B": 2.0,
		"C": 3.0,
		"D": 4.0,
	}
	updated := map[string]float64{
		"A": 5.0,
		"B": 6.0,
		"C": 7.0,
		"D": 8.0,
	}
	s := NewState()
	for name, value := range data {
		s.Add(name, value)
	}
	for name, value := range data {
		vv, ok := s.Get(name)
		if !ok {
			t.Fatalf("variable %v missing", name)
		}
		if vv.(float64) != value {
			t.Fatalf("variable %v should have value %v not %v", name, value, vv)
		}
	}
	for name, value := range updated {
		s.Update(name, value)
	}
	for name, value := range updated {
		vv, ok := s.Get(name)
		if !ok {
			t.Fatalf("variable %v missing", name)
		}
		if vv.(float64) != value {
			t.Fatalf("variable %v should have value %v not %v", name, value, vv)
		}
	}
}

func TestIsSubsetOf(t *testing.T) {
	c1 := NewState()
	c2 := NewState()
	c1.Add("a", 1.0)
	c1.Add("b", -1.0)
	c2.Add("a", 1.0)
	c2.Add("b", -1.0)
	c2.Add("c", 0.0)
	if ok := c1.IsSubsetOf(*c2); !ok {
		t.Error("c1 must be part of c2!")
	}
	if ok := c2.IsSubsetOf(*c1); ok {
		t.Error("c2 cannot be part of c1!")
	}
}

func TestClone(t *testing.T) {
	c1 := NewState()
	c1.Add("a", 1.0)
	c1.Add("b", -1.0)
	c2 := c1.Clone()
	if ok := c1.IsSubsetOf(*c2); !ok {
		t.Error("c1 must be part of c2!")
	}
	if ok := c2.IsSubsetOf(*c1); !ok {
		t.Error("c2 must be part of c1!")
	}
	c2.Update("b", 2.0)
	if ok := c1.IsSubsetOf(*c2); ok {
		t.Error("c1 cannot be part of c2 anymore!")
	}
}