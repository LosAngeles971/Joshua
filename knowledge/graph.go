package knowledge

import (
	"errors"
	"it/losangeles971/joshua/state"
)

/*
An edge is a direct connection between two events,
the source is the "cause" event and the target is the "effect" event.
*/
type Edge struct {
	Cause 	*Event `yaml:"cause"`
	Effect 	*Event `yaml:"effect"`
	Outcome	string
}

//e is influenced by ee if at least one event of ee is also an event of ee
func (e Edge) IsInfluencedBy(ee *Edge) (bool, error) {
	if ok, err := e.Cause.IsInfluencedBy(ee.Cause); ok || err != nil {
		return ok, err
	}
	if ok, err := e.Cause.IsInfluencedBy(ee.Effect); ok || err != nil {
		return ok, err
	}
	if ok, err := e.Effect.IsInfluencedBy(ee.Cause); ok || err != nil {
		return ok, err
	}
	if ok, err := e.Effect.IsInfluencedBy(ee.Effect); ok || err != nil {
		return ok, err
	}
	return false, nil
}

// if cause ran successfully, then try to run then effect
// note: if effect failed to occur, the changed applied by cause still remain
func (e *Edge) Run(input state.State) (string, state.State, error) {
	cause_outcome, cause_output, err := e.Cause.Run(input)
	if err != nil {
		return EFFECT_OUTCOME_ERROR, cause_output, err
	}
	if cause_outcome == EVENT_OUTCOME_FALSE {
		return EFFECT_OUTCOME_CAUSE_FALSE, cause_output, nil
	}
	if cause_outcome == EVENT_OUTCOME_UNKNOWN {
		return EFFECT_OUTCOME_UNKNOWN, cause_output, nil
	}
	effect_outcome, effect_output, err := e.Effect.Run(cause_output)
	if err != nil {
		return EFFECT_OUTCOME_ERROR, effect_output, err
	}
	if effect_outcome == EVENT_OUTCOME_FALSE {
		return EFFECT_OUTCOME_EFFECT_FALSE, effect_output, nil
	}
	if effect_outcome == EVENT_OUTCOME_UNKNOWN {
		return EFFECT_OUTCOME_UNKNOWN, effect_output, nil
	}
	return EFFECT_OUTCOME_TRUE, effect_output, nil
}

/*
A path is the concatenation of edges.
From a cause-effect perspective, a path represents how one "source" event can cause a far "target" event,
thorughout a chain of concatenated cause-effect events.
*/
type Path struct {
	Path 		[]*Edge			`yaml:"path"`
	Executed 	bool			`yaml:"executed"`
	Input 		state.State		`yaml:"input"`
	Output 		state.State		`yaml:"output"`
	Outcome 	string			`yaml:"outcome"`
	Changed 	bool			`yaml:"changed"`
	Cycle 		int				`yaml:"cycle"`
}

// Clone a Path resetting its fields
func (p *Path) clone() *Path {
	n := Path{
		Path: []*Edge{},
		Executed: false,
		Outcome: EFFECT_OUTCOME_NULL,
		Changed: false,
		Cycle: -1,
	}
	for _, e := range p.Path {
		ee := Edge{
			Cause: e.Cause,
			Effect: e.Effect,
			Outcome: EFFECT_OUTCOME_NULL,
		}
		n.Path = append(n.Path, &ee)
	}
	return &n
}

// Used by genetic library
func (p *Path) GetOutcome() string {
	return p.Outcome
}

// Used by genetic library a Path to not executed
func (p *Path) Reset() {
	p.Executed = false
}

func (p *Path) GetWeight() float64 {
	w := float64(0.0)
	for _, e := range p.Path {
		w += e.Cause.GetWeightTo(e.Effect)
	}
	return w
}

func (p *Path) Contains(ee *Event) bool {
	for _, e := range p.Path {
		if e.Cause.GetID() == ee.GetID() || e.Effect.GetID() == ee.GetID() {
			return true
		}
	}
	return false
}

/*
This method checks if given path starts with the receiver path
*/
func (p *Path) IsPrefixOf(pp *Path) bool {
	if len(p.Path) > len(pp.Path) {
		return false
	}
	for i := range p.Path {
		if p.Path[i].Cause.GetID() != pp.Path[i].Cause.GetID() {
			return false
		}
		if p.Path[i].Effect.GetID() != pp.Path[i].Effect.GetID() {
			return false
		}
	}
	return true
}

func (p *Path) equals(pp *Path) bool {
	if p.IsPrefixOf(pp) && pp.IsPrefixOf(p) {
		return true
	}
	return false
}

// This function run a path starting from a specific index
func (p *Path) RunFromTail(input state.State, tail int) (error) {
	for i := tail; i >= 0; i-- {
		e := p.Path[i]
		outcome, output, err := e.Run(input)
		if err != nil {
			e.Outcome = EFFECT_OUTCOME_ERROR
			return err
		}
		e.Outcome = outcome
		p.Outcome = outcome
		p.Output = output
		if !p.Input.IsSubsetOf(p.Output) {
			p.Changed = true
		}
	}
	return nil
}

/* 
Argument input comes from the execution of previous paths
and override the internal input of the receiver's struct.

NOTE: if the first item is the last effect of the chain

The verification must got for attempts:
- just the first rel
- if not the first -1 and then the first
- ...
- if not the first - n, first - n -1, ..., the the first
*/
func (p *Path) Run(input state.State, cycle int) (error) {
	p.Cycle = cycle
	if p.Executed {
		return errors.New("asked to run an already executed state")
	}
	p.Input = input.Clone()
	p.Executed = true
	p.Changed = false
	for i := 0; i < len(p.Path); i++ {
		err := p.RunFromTail(input, i)
		if err != nil {
			return err
		}
		input = p.Output.Clone()
		switch p.Outcome {
		case EFFECT_OUTCOME_CAUSE_FALSE:
			return nil;
		case EFFECT_OUTCOME_EFFECT_FALSE:
			return nil
		case EFFECT_OUTCOME_ERROR:
			return errors.New("outcome is error")
		case EFFECT_OUTCOME_LOOP:
		case EFFECT_OUTCOME_NULL:
			return errors.New("outcome is null")
		case EFFECT_OUTCOME_TRUE:
		case EFFECT_OUTCOME_UNKNOWN:
			return nil
		}
	}
	return nil
}

// create a path between two events
func createPath(cause *Event, effect *Event) Path {
	e := Edge{
		Cause: cause,
		Effect: effect,
	}
	return Path{
		Path: []*Edge{&e},
		Executed: false,
		Changed: false,
		Outcome: EFFECT_OUTCOME_NULL,
		Cycle: -1,
	}
}

func getBranch(o *Path, e *Edge) Path {
	p := Path{
		Path: append(o.Path, e),
		Executed: false,
		Changed: false,
		Outcome: EFFECT_OUTCOME_NULL,
		Cycle: -1,
	}
	return p
}

// a paths is influenced by p if at least one of the s's edge is influenced by one of the p's edge
func (s Path) IsInfluencedBy(p Path) (bool, error) {
	for _, p_e := range p.Path {
		for _, s_e := range s.Path {
			if ok, err := s_e.IsInfluencedBy(p_e); ok || err != nil {
				return ok, err
			}
		}
	}
	return false, nil
}

// Here it is necessary to avoid loops
func backward(p *Path, k Knowledge, s *Stack) {
	last := p.Path[len(p.Path) - 1]
	effect := last.Cause
	causes := k.WhoCause(*effect)
	for i, cause := range causes {
		e := Edge{
			Cause: cause,
			Effect: effect,
		}
		// If the cause is already into the Path, the Path ends here
		// to avoid loops
		if !p.Contains(cause) {
			if i == 0 {
				p.Path = append(p.Path, &e)
				s.Push(p)
			} else {
				b := getBranch(p, &e)
				s.Push(&b)
			}
		}
	}
}