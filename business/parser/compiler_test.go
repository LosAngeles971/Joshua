package parser

import (
	"testing"
)

func TestCompiler(t *testing.T) {
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
	ee, err := Build(parser.code)
	if err != nil {
		t.Fatal(err)
	}
	if len(ee) != len(parser.code) {
		t.Fatalf("expected %v events not %v", len(parser.code), len(ee))
	}
}