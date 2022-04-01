package math

import (
	"testing"
	"fmt"

	"github.com/Knetic/govaluate"
)

// func TestCompleteness(t *testing.T) {
// 	expr, err := govaluate.NewEvaluableExpression("A*(2+3)")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	ok := isComplete(expr, *NewState())
// 	if ok {
// 		t.Fatal("unexpected complete expression")
// 	}
// }

var expr_tests = map[string]interface{}{
	"10+4*5": float64(30),
	"min(1, 2.0, -3, A)": float64(-4.0),
	"A > -10 && A < 0": true,
	"A == -4.0": true,
}

var env_tests = map[string]interface{}{
	"A": -4.0,
}

func TestFunctions(t *testing.T) {
	min, ok := Functions["min"]
	if !ok {
		fmt.Println("Function min does not exist")
		t.FailNow()
	}
	c, err := min(1, 2, 3, 4.0, 5)
	if err != nil {
		fmt.Println("Error: ", err)
		t.FailNow()
	}
	if c != 1.0 {
		fmt.Println("Wrong result: ", c)
		t.FailNow()
	}
}

func TestGovaluate(t *testing.T) {
	for expr, result := range expr_tests {
		eval, err := govaluate.NewEvaluableExpressionWithFunctions(expr, Functions)
		if err != nil {
			t.Fatal(err)
		}
		r, err := eval.Evaluate(env_tests)
		if err != nil {
			t.Fatal(err)
		}
		switch result.(type) {
		case bool:
			if result.(bool) != r.(bool) {
				t.Fatalf("expected %v not %v", result, r)
			}
		default:
			if result.(float64) != r.(float64) {
				t.Fatalf("expected %v not %v", result, r)
			}
		}
	}
}