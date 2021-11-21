package math

import (
	"it/losangeles971/joshua/state"
	"testing"

	"github.com/Knetic/govaluate"
)

func TestCompleteness(t *testing.T) {
	expr, err := govaluate.NewEvaluableExpression("A*(2+3)")
	if err != nil {
		t.Fatal(err)
	}
	ok := IsComplete(expr, state.NewSimpleState())
	if ok {
		t.Fatal("unexpected complete expression")
	}
}