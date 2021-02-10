package knowledge

import (
	"fmt"
	"testing"
	ctx "it/losangeles971/joshua/internal/context"
)

func TestEvents(t *testing.T) {
	init := ctx.Create()
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
		Statements: []string{
			"LupiA + CapreA + CavoliA + ContadiniA == 4",
		},
		Assignments: []string{
			"ContadiniB = 1",
		},
	}
	data := init.Clone()
	outcome, err := e1.Verify(&data)
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