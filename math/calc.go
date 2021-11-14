package math

import (
	"it/losangeles971/joshua/state"
	"log"

	"github.com/Knetic/govaluate"
)

func IsComplete(expr *govaluate.EvaluableExpression, data state.State) bool {
	for _, k := range expr.Vars() {
		vv, ok := data.Get(k)
		if !ok {
			log.Printf("missing variabile %v from state", k)
			return false
		}
		if !vv.IsDefined() {
			log.Printf("variabile %v of state is not defined", k)
			return false
		}
	}
	return true
}