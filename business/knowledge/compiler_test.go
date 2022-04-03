package knowledge

import (
	_ "embed"
	"testing"
	"it/losangeles971/joshua/business/parser"
)

//go:embed thefarmer.joshua
var thefarmer string

func TestCompiler(t *testing.T) {
	scanner, err := parser.NewScanner(thefarmer)
	if err != nil {
		t.Fatal(err)
	}
	lexer, err := scanner.Run()
	if err != nil {
		t.Fatal(err)
	}
	p, err := parser.NewParser(lexer)
	if err != nil {
		t.Fatal(err)
	}
	code, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	ee, err := build(code)
	if err != nil {
		t.Fatal(err)
	}
	if len(ee) != len(code) {
		t.Fatalf("expected %v events not %v", len(code), len(ee))
	}
}