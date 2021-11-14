package knowledge

import (
	"io/ioutil"
	"testing"
)

var	EVENTS = []string{
		"They are all on the est bank of the river",
		"The farmer brings the cabbage to the est bank of the river",
		"The farmer brings the cabbage to the ovest bank of the river",
		"The farmer brings the goat to the est bank of the river",
		"The farmer brings the goat to the ovest bank of the river",
		"The farmer brings the wolf to the est bank of the river",
		"The farmer brings the wolf to the ovest bank of the river",
		"The farmer goes to the est bank of the river",
		"The farmer comes back to the ovest bank of the river",
	}

func TestParse(t *testing.T) {
	source, err := ioutil.ReadFile("../../../resources/the_farmer.joshua")
	if err != nil {
		t.Fatal(err)
	}
	ee, err := Parse(string(source))
	if err != nil {
		t.Fatal(err)
	}
	for i := range EVENTS {
		ok := false
		for j := range ee {
			if ee[j].GetID() == EVENTS[i] {
				ok = true
			}
		}
		if !ok {
			t.Fatalf("missing event %v", EVENTS[i])
		}
	}
}
