package math

import (
	"it/losangeles971/joshua/state"
	"log"

	"github.com/Knetic/govaluate"
)

// Check if the state includes all needed variables by the given expression
func IsComplete(expr *govaluate.EvaluableExpression, data state.State) bool {
	for _, k := range expr.Vars() {
		_, ok := data.Get(k)
		if !ok {
			log.Printf("missing variabile %v from state", k)
			return false
		}
	}
	return true
}