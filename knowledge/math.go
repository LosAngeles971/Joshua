package knowledge

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Knetic/govaluate"
	log "github.com/sirupsen/logrus"
)

// Functions includes common mathematical functions usable into the event's expressions
var Functions = map[string]govaluate.ExpressionFunction{
	"min": func(args ...interface{}) (interface{}, error) {
		min := 0.0
		for i := 0; i < len(args); i ++ {
			switch args[i].(type) {
			case int:
				v := float64(args[i].(int))
				if i == 0 {
					min = v
				} else {
					if v < min {
						min = v
					}
				}
			case float64:
				v := args[i].(float64)
				if i == 0 {
					min = v
				} else {
					if v < min {
						min = v
					}
				}
			default:
				return min, errors.New("not numerical inputs")
			}
		}
		return min, nil
	},
	"max": func(args ...interface{}) (interface{}, error) {
		max := 0.0
		for i := 0; i < len(args); i ++ {
			switch args[i].(type) {
			case int:
				v := float64(args[i].(int))
				if i == 0 {
					max = v
				} else {
					if v > max {
						max = v
					}
				}
			case float64:
				v := args[i].(float64)
				if i == 0 {
					max = v
				} else {
					if v > max {
						max = v
					}
				}
			default:
				return max, errors.New("not numerical inputs")
			}
		}
		return max, nil
	},	
}

// Check if the state includes all needed variables by the given expression
func isComplete(expr *govaluate.EvaluableExpression, data State) bool {
	for _, k := range expr.Vars() {
		_, ok := data.Get(k)
		if !ok {
			log.Printf("missing variabile %v from state", k)
			return false
		}
	}
	return true
}

// ParseExpression returns the object representing a parsed mathematical expression
func parseExpression(e string) (*govaluate.EvaluableExpression, error) {
	expr, err := govaluate.NewEvaluableExpressionWithFunctions(e, Functions)
	if err != nil {
		return nil, err
	}
	return expr, nil
}

// this method splits an assignment into variable and assignement's expression
func splitAssignment(expr string) (string, string, error) {
	parts := strings.Split(expr, "=")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("malformed assignment: %v", expr)
	}
	return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), nil
}

// ParseAssignment runs an assignment of a variable, returning:
// - the target variable
// - the parsed expression
// - the possible error
func parseAssignment(assignemnt string) (string, *govaluate.EvaluableExpression, error) {
	v, e, err := splitAssignment(assignemnt)
	if err != nil {
		return "", nil, err
	}
	expr, err := parseExpression(e)
	return v, expr, err
}