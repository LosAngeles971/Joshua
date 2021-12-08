package engine

import (
	"testing"
)

const (
	EXPECTED_SUCCESS = "They are all on the est bank of the river"
)
var EXPECTED_VARS = map[string]float64{
	"Farmer_location": 0.0,
	"Wolf_location": 0.0,
	"Goat_location": 0.0,
	"Cabbage_location": 0.0,
}

var the_farmer_problem = `---
  success: "They are all on the est bank of the river"
  variables:
    - name: Farmer_location
      value: 0.0
    - name: Wolf_location
      value: 0.0
    - name: Goat_location
      value: 0.0
    - name: Cabbage_location
      value: 0.0
`

func TestLoadProblem(t *testing.T) {
	s, success, err := LoadProblem(the_farmer_problem)
	if err != nil {
		t.Fatal(err)
	}
	if success != EXPECTED_SUCCESS {
		t.Errorf("expected event '%v' not '%v'", EXPECTED_SUCCESS, success)
	}
	for k, v := range EXPECTED_VARS {
		vv, ok := s.Get(k)
		if !ok {
			t.Fatalf("missing variable %v", k)
		}
		if vv.(float64) != v {
			t.Errorf("expected value %v not %v", v, vv)
		}
	}
}


	