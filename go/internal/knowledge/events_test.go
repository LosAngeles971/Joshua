package knowledge

import (
	"fmt"
	"testing"
	ctx "it/losangeles971/joshua/internal/context"
)

/*
This test verifies the functionalities of:
- conditions verification
- assignements execution
*/
func TestEvents(t *testing.T) {
	init := ctx.CreateEmptyState()
	init.Add(&ctx.Variable{Name: "ContadiniA", Value: 1.0, Defined: true, })
	init.Add(&ctx.Variable{Name: "ContadiniB", Value: 0.0, Defined: true, })
	init.Add(&ctx.Variable{Name: "LupiA", Value: 1.0, Defined: true, })
	init.Add(&ctx.Variable{Name: "LupiB", Value: 0.0, Defined: true, })
	init.Add(&ctx.Variable{Name: "CapreA", Value: 1.0, Defined: true, })
	init.Add(&ctx.Variable{Name: "CapreB", Value: 0.0, Defined: true, })
	init.Add(&ctx.Variable{Name: "CavoliA", Value: 1.0, Defined: true, })
	init.Add(&ctx.Variable{Name: "CavoliB", Value: 0.0, Defined: true, })
	e1 := Event{
		ID: "1",
		Conditions: []string{
			"LupiA + CapreA + CavoliA + ContadiniA == 4",
		},
		Assignments: []string{
			"ContadiniB = 1",
		},
	}
	data := init.Clone()
	outcome, err := e1.CanHappen(&data)
	if err != nil {
		fmt.Println("Error ", err)
		t.FailNow()
	}
	if outcome != EVENT_OUTCOME_TRUE {
		fmt.Println("Outcome is not ", EVENT_OUTCOME_TRUE)
		t.FailNow()
	}
	if v, ok := data.Get("ContadiniB"); !ok || v.Value != 1 {
		fmt.Println("ContadiniB is not 1")
		t.FailNow()
	}
}

/*
This test verifies that an event cannot happen if all variables it requires
are defined into the state
*/
func TestUnknownEvents(t *testing.T) {
	init := ctx.CreateEmptyState()
	init.Add(&ctx.Variable{Name: "ContadiniA", Value: 1.0, Defined: true, })
	init.Add(&ctx.Variable{Name: "ContadiniB", Value: 0.0, Defined: true, })
	init.Add(&ctx.Variable{Name: "LupiA", Value: 1.0, Defined: true, })
	init.Add(&ctx.Variable{Name: "LupiB", Value: 0.0, Defined: true, })
	init.Add(&ctx.Variable{Name: "CapreA", Value: 1.0, Defined: true, })
	init.Add(&ctx.Variable{Name: "CapreB", Value: 0.0, Defined: true, })
	init.Add(&ctx.Variable{Name: "CavoliA", Value: 1.0, Defined: false, })
	init.Add(&ctx.Variable{Name: "CavoliB", Value: 0.0, Defined: false, })
	e1 := Event{
		ID: "1",
		Conditions: []string{
			"LupiA + CapreA + CavoliA + ContadiniA == 4",
		},
		Assignments: []string{
			"ContadiniB = 1",
		},
	}
	data := init.Clone()
	outcome, err := e1.CanHappen(&data)
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

/*
Test weight of events
*/