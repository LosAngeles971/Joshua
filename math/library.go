package math

import (
	"errors"
	"github.com/Knetic/govaluate"
)

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