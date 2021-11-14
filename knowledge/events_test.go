package knowledge

import (
	"fmt"
	"testing"
	"it/losangeles971/joshua/state"
)

/*
This test verifies the functionalities of:
- conditions verification
- assignements execution
*/
func TestEvents(t *testing.T) {
	s := state.New()
	s.Add("ContadiniA", 1.0)
	s.Add("LupiA", 1.0)
	s.Add("CapreA", 1.0)
	s.Add("CavoliA", 1.0)
	s.Add("ContadiniB", 0.0)
	e1 := NewEvent("Il contadino va sulla sponda B")
	e1.AddConditions([]string{"ContadiniA + LupiA + CapreA + CavoliA == 4",})
	e1.AddAssignments([]string{"ContadiniB = 1",})
	outcome, output, err := e1.Run(*s)
	if err != nil {
		t.Fatal(err)
	}
	if outcome != EVENT_OUTCOME_TRUE {
		t.Fatalf("outcome should be %v not %v", EVENT_OUTCOME_TRUE, outcome)
	}
	v, ok := output.Get("ContadiniB")
	if !ok {
		t.Error("ContadiniB is not present")
	}
	d, ok := v.GetValue()
	if !ok {
		t.Error("ContadiniB is not defined")
	}
	if d.(float64) != 1.0 {
		t.Errorf("ContadiniB is not 1.0 by %v", d.(float64))
	}
}

/*
This test verifies that an event cannot happen if all variables it requires
are defined into the state
*/
func TestUnknownEvents(t *testing.T) {
	init := state.New()
	init.Add("ContadiniA", 1.0)
	init.Add("ContadiniB", 0.0)
	init.Add("LupiA", 1.0)
	init.Add("LupiB", 0.0)
	init.Add("CapreA", 1.0)
	init.Add("CapreB", 0.0)
	init.Declare("CavoliA", state.TYPE_NUMBER)
	init.Declare("CavoliB", state.TYPE_NUMBER)
	e1 := NewEvent("1")
	e1.AddConditions([]string{"LupiA + CapreA + CavoliA + ContadiniA == 4",})
	e1.AddAssignments([]string{"ContadiniB = 1",})
	data := init.Clone()
	outcome, _, err := e1.Run(data)
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