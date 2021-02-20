package knowledge

import (
	ctx "it/losangeles971/joshua/internal/context"
)

type Edge struct {
	Cause 	*Event	`yaml:"cause"`
	Effect 	*Event	`yaml:"effect"`
}

func (e Edge) IsInfluencedBy(ee Edge) (bool, error) {
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

type Path struct {
	Path 		[]Edge			`yaml:"path"`
	executed 	bool			`yaml:"executed"`
	input 		ctx.State		`yaml:"input"`
	output 		ctx.State		`yaml:"output"`
	outcome 	string			`yaml:"outcome"`
	changed 	bool			`yaml:"changed"`
	cycle 		int				`yaml:"cycle"`
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

func (p *Path) Equals(pp *Path) bool {
	if p.IsPrefixOf(pp) && pp.IsPrefixOf(p) {
		return true
	}
	return false
}

func getPath(cause *Event, effect *Event) Path {
	e := Edge{
		Cause: cause,
		Effect: effect,
	}
	return Path{
		Path: []Edge{e},
		executed: false,
		changed: false,
		outcome: CE_OUTCOME_NULL,
		cycle: -1,
	}
}

func getBranch(o *Path, e Edge) Path {
	p := Path{
		Path: append(o.Path, e),
		executed: false,
		changed: false,
		outcome: CE_OUTCOME_NULL,
		cycle: -1,
	}
	return p
}

// Here it is necessary to avoid loops
func backward(p *Path, k Knowledge, s *Stack) {
	last := p.Path[len(p.Path) - 1]
	effect := last.Cause
	causes := k.IsEffectOf(effect)
	for i, cause := range causes {
		e := Edge{
			Cause: cause,
			Effect: effect,
		}
		// If the cause is already into the Path, the Path ends here
		// to avoid loops
		if !p.Contains(cause) {
			if i == 0 {
				p.Path = append(p.Path, e)
				s.Push(p)
			} else {
				b := getBranch(p, e)
				s.Push(&b)
			}
		}
	}
} 

/*
This method finds all paths into the Knowledge
that produce as final effect the given effect.
*/
func GetAllPaths(k Knowledge, effect *Event) []*Path {
	s := &Stack{}
	discovered := []*Path{}
	for _, cause := range k.IsEffectOf(effect) {
		p := getPath(cause, effect)
		discovered = append(discovered, &p)
		s.Push(&p)
	}
	for s.Size() > 0 {
		p, ok := s.Pop()
		if!ok {
			return discovered
		}
		backward(p, k, s)
	}
	return discovered
}