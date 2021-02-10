package runtime

import (
	"errors"
	ctx "it/losangeles971/joshua/internal/context"
	kkk "it/losangeles971/joshua/internal/knowledge"
)

type Path struct {
	rel kkk.Relationship
	executed bool
	input ctx.State
	output ctx.State
	outcome string
	changed bool
	cycle int
}

// ATTENTION: the input context must be updated due to the actions of the
// previous executed states
func (s *Path) Run(input ctx.State, cycle int) (error) {
	s.cycle = cycle
	if s.executed {
		return errors.New("Asked to run an already executed state")
	}
	s.input = input.Clone()
	s.executed = true
	s.changed = false
	var err error
	outcome, output, err := s.rel.Verify(s.input)
	if err != nil {
		return err
	}
	s.outcome = outcome
	s.output = output.Clone()
	if !s.input.PartOf(s.output) {
		s.changed = true
	}
	return nil
}

func (s Path) Outcome() string {
	return s.outcome
}

func (s Path) Changed() bool {
	return s.changed
}

func (s Path) Cycle() int {
	return s.cycle
}

func (s Path) Executed() bool {
	return s.executed
}

func (s Path) Input() ctx.State {
	return s.input.Clone()
}

func (s Path) Output() ctx.State {
	return s.output.Clone()
}

func (s *Path) Loop() {
	s.outcome = kkk.CE_OUTCOME_LOOP
}

func (s Path) Relationship() kkk.Relationship {
	return s.rel
}

func (s Path) IsInfluenced(rel kkk.Relationship) (bool, error) {
	ok, err := s.rel.IsInfluencedBy(rel)
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}
	return false, nil
}