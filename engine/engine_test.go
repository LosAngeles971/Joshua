package engine

import (
	"io/ioutil"
	"it/losangeles971/joshua/knowledge"
	"it/losangeles971/joshua/state"
	"log"
	"testing"
)

func TestLogicReasoning(t *testing.T) {
	source, err := ioutil.ReadFile("../../../resources/the_farmer.joshua")
	if err != nil {
		t.Fatal(err)
	}
	k, err := knowledge.Load(string(source))
	if err != nil {
		t.Fatal(err)
	}
	s := state.New()
	s.Add("Farmer_location", 0.0)
	s.Add("Wolf_location", 0.0)
	s.Add("Goat_location", 0.0)
	s.Add("Cabbage_location", 0.0)
	success, ok := k.GetEvent("They are all on the est bank of the river")
	if !ok {
		t.Fatal("missing success event")
	}
	outcome, solution, err := MakeItHappen(k, *s, success, 100)
	if err != nil {
		t.Fatal(err)
	}
	if solution.Size() < 1 {
		t.Fatal("expected a queue larger than 0")
	}
	if outcome != knowledge.EFFECT_OUTCOME_TRUE {
		t.Fatalf("exepected outcome [%v] not [%v]", knowledge.EFFECT_OUTCOME_TRUE, outcome)
	}
	PrintSummary(outcome, solution)
	y, err := solution.Serialize(true)
	if err != nil {
		t.Fatal(err)
	}
	log.Print(y)
}


	