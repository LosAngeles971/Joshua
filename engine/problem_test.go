package engine

import (
	"io/ioutil"
	"testing"
)

const (
	EXPECTED_SUCCES = "They are all on the est bank of the river"
)
var EXPECTED_VARS = map[string]float64{
	"Farmer_location": 0.0,
	"Wolf_location": 0.0,
	"Goat_location": 0.0,
	"Cabbage_location": 0.0,
}

func TestLoadProblem(t *testing.T) {
	source, err := ioutil.ReadFile("../../../resources/the_farmer.yml")
	if err != nil {
		t.Fatal(err)
	}
	s, success, err := LoadProblem(string(source))
	if err != nil {
		t.Fatal(err)
	}
	if success != EXPECTED_SUCCES {
		t.Errorf("expected event '%v' not '%v'", EXPECTED_SUCCES, success)
	}
	for k, v := range EXPECTED_VARS {
		vv, ok := s.Get(k)
		if !ok {
			t.Fatalf("missing variable %v", k)
		}
		ff, ok := vv.GetValue()
		if !ok {
			t.Errorf("variable %v should be defined", k)
		}
		if ff.(float64) != v {
			t.Errorf("expected value %v not %v", v, ff)
		}
	}
}


	