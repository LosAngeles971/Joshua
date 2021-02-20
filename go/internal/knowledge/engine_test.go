package knowledge

import (
	"fmt"
	"testing"
	ctx "it/losangeles971/joshua/internal/context"
)

func TestLogicReasoning(t *testing.T) {
	k := Knowledge{}
	err := k.Load("../../../resources/k_contadino.yml")
	if err != nil {
		fmt.Println("Knowledge not loaded due to error ", err)
		t.FailNow()
	}
	init := ctx.CreateEmptyState()
	init.Add(&ctx.Variable{Name: "FarmerA", Value: 1.0, Defined: true, })
	init.Add(&ctx.Variable{Name: "FarmerB", Value: 0.0, Defined: true, })
	init.Add(&ctx.Variable{Name: "WolfA", Value: 1.0, Defined: true, })
	init.Add(&ctx.Variable{Name: "WolfB", Value: 0.0, Defined: true, })
	init.Add(&ctx.Variable{Name: "GoatA", Value: 1.0, Defined: true, })
	init.Add(&ctx.Variable{Name: "GoatB", Value: 0.0, Defined: true, })
	init.Add(&ctx.Variable{Name: "CabbageA", Value: 1.0, Defined: true, })
	init.Add(&ctx.Variable{Name: "CabbageB", Value: 0.0, Defined: true, })
	c_name := "The farmer brings the goat on the bank B of the river"
	e_name := "The farmer, the wolf, the goat and the cabbage are on the bank B of the river"
	_, ok := k.GetEvent(c_name)
	if !ok {
		fmt.Println("Knowledge does not contain event: ", c_name)
		t.FailNow()
	}
	_, ok = k.GetEvent(e_name)
	if !ok {
		fmt.Println("Knowledge does not contain event: ", e_name)
		t.FailNow()
	}
}

func TestCloningPaths(t *testing.T) {
	k := Knowledge{}
	err := k.Load("../../../resources/k_contadino.yml")
	if err != nil {
		fmt.Println("Knowledge not loaded due to error ", err)
		t.FailNow()
	}
	v1, _ := k.GetEvent("The farmer brings the cabbage on the bank B of the river")
	v2, _ := k.GetEvent("The farmer brings the cabbage on the bank A of the river")
	success, _ := k.GetEvent("The farmer, the wolf, the goat and the cabbage are on the bank B of the river")
	e1 := Edge{ Cause: v1, Effect: success, }
	e2 := Edge{ Cause: v2, Effect: success,	}
	p1 := Path{
		Path: []Edge{
			e1,
			e2,
		},
	}
	init := ctx.CreateEmptyState()
	q := Queue{[]*Path{&p1}}
	p1.Run(init, 0)
	p1.outcome = "Test"
	if !p1.executed {
		fmt.Println("P1 should be executed")
		t.FailNow()
	}
	q.addClone(&p1)
	if q.Size() != 2 {
		fmt.Println("Wrong queue size: ", q.Size())
		t.FailNow()
	}
	p2 := q.paths[1]
	if p2.executed {
		fmt.Println("P1 should NOT be executed")
		t.FailNow()
	}
	p2.outcome = "Weird"
	if p1.outcome == p2.outcome {
		fmt.Println("Entanglement!!!")
		t.FailNow()
	}
	if p2.outcome != "Weird" || p1.outcome != "Test" {
		fmt.Println("Data corruption")
		t.FailNow()
	}
}
	