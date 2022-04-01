package knowledge

import (
	"it/losangeles971/joshua/business/math"

	"github.com/Knetic/govaluate"
)

// Assignment sets a value to a variabile, as a consequence of the owner event's occurrence
type Assignment struct {
	variable string
	expr     *govaluate.EvaluableExpression
}

func NewAssignment(expr string) (Assignment, error) {
	v, e, err := math.NewAssignment(expr)
	if err != nil {
		return Assignment{}, err
	}
	return Assignment{
		variable: v,
		expr: e,
	}, nil
}
