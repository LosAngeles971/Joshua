package knowledge

import (
	_ "embed"
	"testing"
)

//go:embed thefarmer.yml
var test_data1 string
//go:embed thefarmer.json
var test_data2 string

func TestLoad(t *testing.T) {
	s1 := NewState(WithData([]byte(test_data1), ENC_YAML))
	s2 := NewState(WithData([]byte(test_data2), ENC_JSON))
	expected := map[string]int{
		"Farmer_location": 0,
  		"Wolf_location": 0,
  		"Goat_location": 0,
  		"Cabbage_location": 0,
	}
	for name, value := range expected {
		v1, ok1 := s1.Get(name)
		v2, ok2 := s2.Get(name)
		if !ok1 {
			t.Errorf("s1: variable %s is missing", name)
		}
		if !ok2 {
			t.Errorf("s2: variable %s is missing", name)
		}
		if ok1 && v1.(int) != value {
			t.Errorf("s1: variable %s got %v not %v", name, v1, value)
		}
		if ok2 && v2.(float64) != float64(value) {
			t.Errorf("s2: variable %s got %v not %v", name, v2, value)
		}
	}
}

func TestCreationAndUpdate(t *testing.T) {
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