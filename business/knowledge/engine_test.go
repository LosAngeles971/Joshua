package knowledge

import (
	log "github.com/sirupsen/logrus"
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
	outcome, solution, err := engine.IsItGoingToHappen(*s, "They are all on the est bank of the river")
	if err != nil {
		t.Fatal(err)
	}
	if solution.Size() < 1 {
		t.Fatal("expected a queue larger than 0")
	}
	if outcome != EFFECT_OUTCOME_TRUE {
		t.Fatalf("exepected outcome [%v] not [%v]", EFFECT_OUTCOME_TRUE, outcome)
	}
	PrintSummary(outcome, solution)
	y, err := solution.Serialize()
	if err != nil {
		t.Fatal(err)
	}
	log.Print(y)
}


	