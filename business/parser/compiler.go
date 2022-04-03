package parser

import (
	"fmt"
	"it/losangeles971/joshua/business/knowledge"
	"it/losangeles971/joshua/business/math"
	"strconv"

	"github.com/Knetic/govaluate"
)

type EventCode struct {
	name        string
	conditions  []string
	assignments []string
	effects     []string
}

func (code EventCode) getConditions() ([]*govaluate.EvaluableExpression, error) {
	cc := []*govaluate.EvaluableExpression{}
	for i := range code.conditions {
		c, err := math.NewExpression(code.conditions[i])
		if err != nil {
			return nil, err
		}
		cc = append(cc, c)
	}
	return cc, nil
}

func (code EventCode) getAssignments() ([]knowledge.Assignment, error) {
	aa := []knowledge.Assignment{}
	for i := range code.assignments {
		a, err := knowledge.NewAssignment(code.assignments[i])
		if err != nil {
			return nil, err
		}
		aa = append(aa, a)
	}
	return aa, nil
}

func (code EventCode) getRelationships() ([]*knowledge.Relationship, error) {
	rr := []*knowledge.Relationship{}
	i := 0
	for i < len(code.effects) {
		name := code.effects[i]
		i++
		if i >= len(code.effects) {
			return nil, fmt.Errorf("REL [%v][%s] missing weight", i, name)
		}
		weight, err := strconv.ParseFloat(code.effects[i], 64)
		if err != nil {
			return nil, fmt.Errorf("REL [%v][%s] weight is not a float64 number [%s]", i, name, code.effects[i])
		}
		i++
		r := knowledge.NewRelationship(name, knowledge.WithWeight(weight))
		rr = append(rr, r)
	}
	return rr, nil
}

func (code EventCode) compileEvent() (*knowledge.Event, error) {
	cc, err := code.getConditions()
	if err != nil {
		return nil, err
	}
	aa, err := code.getAssignments()
	if err != nil {
		return nil, err
	}
	rr, err := code.getRelationships()
	if err != nil {
		return nil, err
	}
	return knowledge.NewEvent(code.name, 
		knowledge.WithConditions(cc),
		knowledge.WithAssignments(aa),
		knowledge.WithRelationships(rr)), nil
}

func Build(code []EventCode) ([]*knowledge.Event, error) {
	built := []*knowledge.Event{}
	for i := range code {
		e, err := code[i].compileEvent()
		if err != nil {
			return nil, err
		}
		built = append(built, e)
	}
	return built, nil
}

func Compile(source string) ([]*knowledge.Event, error) {
	scanner, err := NewScanner(source)
	if err != nil {
		return nil, err
	}
	lexer, err := scanner.run()
	if err != nil {
		return nil, err
	}
	parser, err := NewParser(lexer)
	if err != nil {
		return nil, err
	}
	err = parser.Parse()
	if err != nil {
		return nil, err
	}
	return Build(parser.code)	
}