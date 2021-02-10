package runtime

import (
	"fmt"
	"testing"
	"it/losangeles971/joshua/internal/io"
	ctx "it/losangeles971/joshua/internal/context"
	kkk "it/losangeles971/joshua/internal/knowledge"
)

func TestLogicReasoning(t *testing.T) {
	k, err := io.Load("../../../resources/k_contadino.yml")
	if err != nil {
		fmt.Println("Knowledge not loaded due to error ", err)
		t.FailNow()
	}
	init := ctx.Create()
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
	cause, ok := k.GetEvent(c_name)
	if !ok {
		fmt.Println("Knowledge does not contain event: ", c_name)
		t.FailNow()
	}
	success, ok := k.GetEvent(e_name)
	if !ok {
		fmt.Println("Knowledge does not contain event: ", e_name)
		t.FailNow()
	}
	queue := Queue{}
	queue.Populate(init, k, success)
	if queue.Size() != 8 {
		fmt.Println("Unexpected queue's size: ", queue.Size())
		t.FailNow()
	}
	rel, ok := k.GetRelationship(cause, success)
	if !ok {
		fmt.Println("Knowledge does not contain desired realationship")
		t.FailNow()
	}
	state := queue.FindByRelationship(rel)
	if state == nil {
		fmt.Println("Chosen state is nil")
		t.FailNow()
	}
	err = state.Run(init, 0)
	if err != nil {
		fmt.Println("State's execution failed due to error ", err)
		t.FailNow()
	}
	if !state.executed || state.cycle != 0 || state.outcome == kkk.CE_OUTCOME_NULL {
		fmt.Println("State's execution dit not change the state's condition")
		t.FailNow()
	}
	if !state.changed {
		fmt.Println("State's execution dit not change the context")
		t.FailNow()
	}
	co_a, _ := state.output.Get("FarmerA")
	co_b, _ := state.output.Get("FarmerB")
	ca_a, _ := state.output.Get("GoatA")
	ca_b, _ := state.output.Get("GoatB")
	if co_a.Value != 0.0 || co_b.Value != 1.0 || ca_a.Value != 0.0 || ca_b.Value != 1.0 {
		fmt.Println("State's context is wrong")
		t.FailNow()
	}
	clone := queue.AddClone(state)
	if clone.executed || clone.cycle != -1 || clone.outcome != kkk.CE_OUTCOME_NULL {
		fmt.Println("Malformed clone")
		t.FailNow()
	}
	if queue.Size() != 9 {
		fmt.Println("Unexpected queue's size: ", queue.Size())
		t.FailNow()
	}
}