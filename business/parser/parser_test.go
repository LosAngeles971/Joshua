package parser

import (
	_ "embed"
	"testing"
)

//go:embed thefarmer.joshua
var test2 string

func TestParser(t *testing.T) {
	events := []string{
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
	scanner, err := NewScanner(test2)
	if err != nil {
		t.Fatal(err)
	}
	lexer, err := scanner.Run()
	if err != nil {
		t.Fatal(err)
	}
	p, err := NewParser(lexer)
	if err != nil {
		t.Fatal(err)
	}
	code, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	if len(code) != len(events) {
		t.Fatalf("expected %v events not %v", len(events), len(code))
	}
	for i := range events {
		if code[i].name != events[i] {
			t.Errorf("expected [%s] not [%s]", events[i], code[i].name)
		}
	}
}