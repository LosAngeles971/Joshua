package math

import (
	"fmt"
	"strings"

	"github.com/Knetic/govaluate"
)

// it returns the parsed expression of a string containing a mathematical expression
func ParseExpression(e string) (*govaluate.EvaluableExpression, error) {
	expr, err := govaluate.NewEvaluableExpressionWithFunctions(e, Functions)
	if err != nil {
		return nil, err
	}
	return expr, nil
}

// this method splits an assignment into variable and assignement's expression
func parseAssignment(expr string) (string, string, error) {
	parts := strings.Split(expr, "=")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("malformed assignment: %v", expr)
	}
	return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), nil
}

// given an assignement, this method returns:
// - the target variable of the assigment
// - the parsed assignment's expression
// - the possible error
func ParseAssignment(assignemnt string) (string, *govaluate.EvaluableExpression, error) {
	v, e, err := parseAssignment(assignemnt)
	if err != nil {
		return "", nil, err
	}
	expr, err := ParseExpression(e)
	return v, expr, err
}