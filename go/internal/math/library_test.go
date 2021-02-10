package math

import (
	"testing"
	"fmt"
	"github.com/Knetic/govaluate"
)

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
	e := "min(1, 2.0, -3, A)"
	eval, err := govaluate.NewEvaluableExpressionWithFunctions(e, Functions)
	if err != nil {
		fmt.Println("Error: ", err)
		t.FailNow()
	}
	vars := map[string]interface{}{
		"A": -4.0,
	}
	r, err := eval.Evaluate(vars)
	if err != nil {
		fmt.Println("Error: ", err)
		t.FailNow()
	}
	v := r.(float64)
	if v != -4.0 {
		fmt.Println("Wrong result: ", v)
		t.FailNow()
	}
}