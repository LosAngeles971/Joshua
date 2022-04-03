package knowledge

import (
	"testing"
)

func TestLogicReasoning(t *testing.T) {
	s := NewState()
	s.Add("Farmer_location", 0.0)
	s.Add("Wolf_location", 0.0)
	s.Add("Goat_location", 0.0)
	s.Add("Cabbage_location", 0.0)
	engine, err := NewEngine(thefarmer)
	if err != nil {
		t.Fatal(err)
	}
	solution := engine.IsItGoingToHappen(*s, "They are all on the est bank of the river")
	if solution.Err != nil {
		t.Fatal(err)
	}
	if solution.Chain.Size() < 1 {
		t.Fatal("expected a queue larger than 0")
	}
	if solution.Outcome != EFFECT_OUTCOME_TRUE {
		t.Fatalf("exepected outcome [%v] not [%v]", EFFECT_OUTCOME_TRUE, solution.Outcome)
	}
	solution.PrintChain()
	solution.PrintSummary()
}


	