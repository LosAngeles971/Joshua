package math

import (
	"fmt"
	"strings"

	"github.com/Knetic/govaluate"
)

func parseAssignment(expr string) (string, string, error) {
	parts := strings.Split(expr, "=")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("malformed assignment: %v", expr)
	}
	return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), nil
}

func ParseAssignment(assignemnt string) (string, *govaluate.EvaluableExpression, error) {
	v, e, err := parseAssignment(assignemnt)
	if err != nil {
		return "", nil, err
	}
	expr, err := govaluate.NewEvaluableExpressionWithFunctions(e, Functions)
	if err != nil {
		return "", nil, err
	}
	return v, expr, nil
}

func ParseExpression(e string) (*govaluate.EvaluableExpression, error) {
	expr, err := govaluate.NewEvaluableExpressionWithFunctions(e, Functions)
	if err != nil {
		return nil, err
	}
	return expr, nil
}