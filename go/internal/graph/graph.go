package internal

/*

import (
	"errors"
	"fmt"
	kkk "it/losangeles971/joshua/internal/knowledge"
	ctx "it/losangeles971/joshua/internal/context"
)

type Path struct {
	Steps 	[]kkk.Relationship
	Enabled bool
}

func (p Path) Print() string {
	output := ""
	for _, r := range p.Steps {
		output += r.Print() + " \n"
	}
	return output
}

func (p Path) Weight() float64 {
	var w float64
	for i, r := range p.Steps {
		if i == 0 {
			w = r.Weight()
		} else {
			w *= r.Weight()
		}
	}
	return w
}

func (p Path) IsTrueNow() (int, error) {
	for _, r := range p.Steps {
		outcome, err := r.IsTrueNow()
		if err != nil {
			return CE_OUTCOME_ERROR, err
		}
		if outcome == CE_OUTCOME_CAUSE_FALSE || outcome == CE_OUTCOME_EFFECT_FALSE || outcome == CE_OUTCOME_UNKNOWN {
			return outcome, nil
		}
	}
	return CE_OUTCOME_TRUE, nil
}

func (p Path) Last() (Relationship, bool) {
	if len(p.Steps) > 0 {
		return p.Steps[len(p.Steps) - 1], true
	}
	return Relationship{}, false
}

func (p Path) First() (Relationship, bool) {
	if len(p.Steps) > 0 {
		return p.Steps[0], true
	}
	return Relationship{}, false
}

func (p Path) Contains(f Event) bool {
	for _, r := range p.Steps {
		if r.Cause.GetID() == f.GetID() || r.Effect.GetID() == f.GetID() {
			return true
		}
	}
	return false
}

*/