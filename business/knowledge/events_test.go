package knowledge

import (
	_ "embed"
	"fmt"
	"testing"
)

//go:embed testevents.joshua
var test_events string

/*
This test verifies the functionalities of:
- conditions verification
- assignements execution
*/
func TestEvents(t *testing.T) {
	kkk, err := NewKnowledge(WithSource(test_events))
	if err != nil {
		t.Fatal(err)
	}
	s := NewState()
	s.Add("ContadiniA", 1.0)
	s.Add("LupiA", 1.0)
	s.Add("CapreA", 1.0)
	s.Add("CavoliA", 1.0)
	s.Add("ContadiniB", 0.0)
	e1, ok := kkk.GetEvent("Il contadino va sulla sponda B")
	if !ok {
		t.Fatalf("not found event")
	}
	outcome, output, err := e1.Run(*s)
	if err != nil {
		t.Fatal(err)
	}
	if outcome != EVENT_OUTCOME_TRUE {
		t.Fatalf("outcome should be %v not %v", EVENT_OUTCOME_TRUE, outcome)
	}
	value, ok := output.Get("ContadiniB")
	if !ok {
		t.Error("ContadiniB is not present")
	}
	if value.(float64) != 1.0 {
		t.Errorf("ContadiniB is not 1.0 by %v", value.(float64))
	}
}

// This test verifies that an event fails if state misses some variables
func TestUnknownEvents(t *testing.T) {
	kkk, err := NewKnowledge(WithSource(test_events))
	if err != nil {
		t.Fatal(err)
	}
	s := NewState()
	s.Add("ContadiniA", 1.0)
	s.Add("LupiA", 1.0)
	s.Add("CapreA", 1.0)
	e1, ok := kkk.GetEvent("Il contadino va sulla sponda B")
	if !ok {
		t.Fatalf("not found event")
	}
	outcome, _, err := e1.Run(*s)
	if err != nil {
		fmt.Println("Error ", err)
		fmt.Println("Outcome ", outcome)
		t.FailNow()
	}
	if outcome != EVENT_OUTCOME_UNKNOWN {
		fmt.Println("Outcome is not ", EVENT_OUTCOME_UNKNOWN, " but ", outcome)
		t.FailNow()
	}
}