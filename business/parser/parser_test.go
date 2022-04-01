package parser

import (
	_ "embed"
	"testing"
)

//go:embed thefarmer.joshua
var test_source string

func TestParse(t *testing.T) {
	t.Log(test_source)
	s := test_source
	p := NewParser(s)
	_, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
}
