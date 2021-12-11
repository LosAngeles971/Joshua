package knowledge

import (
	"fmt"
	"testing"
)

/*
This test verifies the functionalities of:
- conditions verification
- assignements execution
*/
func TestEvents(t *testing.T) {
	s := NewState()
	s.Add("ContadiniA", 1.0)
	s.Add("LupiA", 1.0)
	s.Add("CapreA", 1.0)
	s.Add("CavoliA", 1.0)
	s.Add("ContadiniB", 0.0)
	e1 := newEvent("Il contadino va sulla sponda B")
	e1.addConditions([]string{"ContadiniA + LupiA + CapreA + CavoliA == 4",})
	e1.addAssignments([]string{"ContadiniB = 1",})
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

/*
This test verifies that an event cannot happen if all variables it requires
are defined into the state
*/
func TestUnknownEvents(t *testing.T) {
	init := NewState()
	init.Add("ContadiniA", 1.0)
	init.Add("ContadiniB", 0.0)
	init.Add("LupiA", 1.0)
	init.Add("LupiB", 0.0)
	init.Add("CapreA", 1.0)
	init.Add("CapreB", 0.0)
	e1 := newEvent("1")
	e1.addConditions([]string{"LupiA + CapreA + CavoliA + ContadiniA == 4",})
	e1.addAssignments([]string{"ContadiniB = 1",})
	data := init.Clone()
	outcome, _, err := e1.Run(*data)
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