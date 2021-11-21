package engine

import (
	"it/losangeles971/joshua/knowledge"
	"it/losangeles971/joshua/state"
	"log"
	"testing"
)

var the_farmer =` 
/*
The farmer, the wolf, the goat and the cabbage
*/

event(They are all on the est bank of the river) {
	premises {
	}
	if {source
	  Farmer_location == 1;
	  Wolf_location == 1;
	  Goat_location == 1;
	  Cabbage_location == 1;
	}
	then {
	}
  }
  
  event(The farmer brings the cabbage to the est bank of the river) {
	premises {
	}
	if {
	  Farmer_location == 0;
	  Cabbage_location == 0;
	  Wolf_location != Goat_location;
	}
	then {
	  event(They are all on the est bank of the river, 0.5);
	  Farmer_location = 1;
	  Cabbage_location = 1;
	}
  }
  
  event(The farmer brings the cabbage to the ovest bank of the river) {
	premises {
	}
	if {
	  Farmer_location == 1;
	  Cabbage_location == 1;
	  Wolf_location != Goat_location;
	}
	then {
	  event(They are all on the est bank of the river, 0.1);
	  Farmer_location = 0;
	  Cabbage_location = 0;
	}
  }
  
  event(The farmer brings the goat to the est bank of the river) {
	premises {
	}
	if {
		Farmer_location == 0;
		Goat_location == 0;
	}
	then {
	  event(They are all on the est bank of the river, 0.5);
	  Farmer_location = 1;
	  Goat_location = 1;
	}
  }
  
  event(The farmer brings the goat to the ovest bank of the river) {
	premises {
	}
	if {
		Farmer_location == 1;
		Goat_location == 1;
	}
	then {
	  event(They are all on the est bank of the river, 0.1);
	  Farmer_location = 0;
	  Goat_location = 0;
	}
  }
  
  event(The farmer brings the wolf to the est bank of the river) {
	premises {
	}
	if {
		Farmer_location == 0;
		Wolf_location == 0;
		Cabbage_location != Goat_location;
	}
	then {
	  event(They are all on the est bank of the river, 0.5);
	  Farmer_location = 1;
	  Wolf_location = 1;
	}
  }
  
  event(The farmer brings the wolf to the ovest bank of the river) {
	premises {
	}
	if {
		Farmer_location == 1;
		Goat_location == 1;
		Cabbage_location != Goat_location;
	}
	then {
	  event(They are all on the est bank of the river, 0.1);
	  Farmer_location = 0;
	  Wolf_location = 0;
	}
  }
  
  event(The farmer goes to the est bank of the river) {
	premises {
	}
	if {
	  Farmer_location == 0;
	  (Wolf_location == 1 && Cabbage_location == 1 && Goat_location == 0) ||
	  (Wolf_location == 0 && Cabbage_location == 0 && Goat_location == 1);
	}
	then {
	  event(They are all on the est bank of the river, 0.3);
	  Farmer_location = 1;
	}
  }
  
  event(The farmer comes back to the ovest bank of the river) {
	premises {
	}
	if {
	  Farmer_location == 1;
	  (Wolf_location == 1 && Cabbage_location == 1 && Goat_location == 0) ||
	  (Wolf_location == 0 && Cabbage_location == 0 && Goat_location == 1);
	}
	then {
	  event(They are all on the est bank of the river, 0.3);
	  Farmer_location = 0;
	}
  }
  `

func TestLogicReasoning(t *testing.T) {
	k, err := knowledge.Load(the_farmer)
	if err != nil {
		t.Fatal(err)
	}
	s := state.NewSimpleState()
	s.Add("Farmer_location", 0.0)
	s.Add("Wolf_location", 0.0)
	s.Add("Goat_location", 0.0)
	s.Add("Cabbage_location", 0.0)
	success, ok := k.GetEvent("They are all on the est bank of the river")
	if !ok {
		t.Fatal("missing success event")
	}
	outcome, solution, err := MakeItHappen(k, s, success, 100)
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


	