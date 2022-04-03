package parser

import (
	_ "embed"
	"testing"
)

//go:embed test2.joshua
var test2 string

func TestParser(t *testing.T) {
	scanner, err := NewScanner(test2)
	if err != nil {
		t.Fatal(err)
	}
	lexer, err := scanner.run()
	if err != nil {
		t.Fatal(err)
	}
	parser, err := NewParser(lexer)
	if err != nil {
		t.Fatal(err)
	}
	err = parser.Parse()
	if err != nil {
		t.Fatal(err)
	}
	if len(parser.code) != 2 {
		t.Errorf("expected 2 events not %v", len(parser.code))
	}
	if parser.code[0].ID != "They are all on the est bank of the river" {
		t.Errorf("wrong")
	}
	if parser.code[1].ID != "The farmer brings the cabbage to the est bank of the river" {
		t.Errorf("wrong")
	}
}