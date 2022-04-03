package knowledge

import (
	"fmt"
	"it/losangeles971/joshua/business/math"
	"it/losangeles971/joshua/business/parser"
	"strconv"
	"strings"

	"github.com/Knetic/govaluate"
)

func getConditions(code *parser.EventCode) ([]*govaluate.EvaluableExpression, error) {
	cc := []*govaluate.EvaluableExpression{}
	for _, text := range code.GetConditions() {
		c, err := math.NewExpression(text)
		if err != nil {
			return nil, err
		}
		cc = append(cc, c)
	}
	if len(cc) == 0 {
		return nil, fmt.Errorf("event [%s] has no conditions", code.Name())
	}
	return cc, nil
}

func getAssignments(code *parser.EventCode) ([]Assignment, error) {
	aa := []Assignment{}
	for _, text := range code.GetAssignments() {
		a, err := NewAssignment(text)
		if err != nil {
			return nil, err
		}
		aa = append(aa, a)
	}
	return aa, nil
}

func getRelationships(code *parser.EventCode) ([]*Relationship, error) {
	rr := []*Relationship{}
	i := 0
	ee := code.GetEffects()
	for i < len(ee) {
		name := ee[i]
		i++
		if i >= len(ee) {
			return nil, fmt.Errorf("REL [%v][%s] missing weight", i, name)
		}
		weight, err := strconv.ParseFloat(strings.TrimSpace(ee[i]), 64)
		if err != nil {
			return nil, fmt.Errorf("REL [%v][%s] weight is not a float64 number [%s]", i, name, ee[i])
		}
		i++
		r := NewRelationship(name, WithWeight(weight))
		rr = append(rr, r)
	}
	return rr, nil
}

func compileEvent(code *parser.EventCode) (*Event, error) {
	cc, err := getConditions(code)
	if err != nil {
		return nil, err
	}
	aa, err := getAssignments(code)
	if err != nil {
		return nil, err
	}
	rr, err := getRelationships(code)
	if err != nil {
		return nil, err
	}
	return NewEvent(code.Name(), 
		WithConditions(cc),
		WithAssignments(aa),
		WithRelationships(rr)), nil
}

func build(code []*parser.EventCode) ([]*Event, error) {
	built := []*Event{}
	for i := range code {
		e, err := compileEvent(code[i])
		if err != nil {
			return nil, err
		}
		built = append(built, e)
	}
	return built, nil
}

func link(ee []*Event) error {
	for _, e := range ee {
		for _, effect := range e.effects {
			ok := false
			for _, target := range ee {
				if target.ID == effect.Name {
					effect.Effect = target
					ok = true
				}
			}
			if !ok {
				return fmt.Errorf("effect %v of event %v does not exist", effect.Name, e.ID)
			}
		}
	}
	return nil
}

func Compile(source string) ([]*Event, error) {
	scanner, err := parser.NewScanner(source)
	if err != nil {
		return nil, err
	}
	lexer, err := scanner.Run()
	if err != nil {
		return nil, err
	}
	p, err := parser.NewParser(lexer)
	if err != nil {
		return nil, err
	}
	code, err := p.Parse()
	if err != nil {
		return nil, err
	}
	ee, err := build(code)
	if err != nil {
		return nil, err
	}
	return ee, link(ee)
}