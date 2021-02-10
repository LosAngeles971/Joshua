package runtime

import (
	"errors"
	"fmt"
	"sort"
	kkk "it/losangeles971/joshua/internal/knowledge"
	ctx "it/losangeles971/joshua/internal/context"
)

type Queue struct {
	paths []*Path
}

func (q *Queue) Populate(data ctx.State, k kkk.Knowledge, effect kkk.Event) {
	for _, r := range k.IsEffectOf(effect) {
		s := Path{
			rel: r,
			executed: false,
			input: data,
			changed: false,
			outcome: kkk.CE_OUTCOME_NULL,
			cycle: -1,
		}
		q.paths = append(q.paths, &s)
	}
}

func (q Queue) Size() (int) {
	return len(q.paths)
}

func (q Queue) FindByRelationship(rel kkk.Relationship) (*Path) {
	for _, s := range q.paths {
		if s.rel.Equals(rel) {
			return s
		}
	}
	return nil
}

func (q Queue) Choose() (*Path) {
	if len(q.paths) < 1 {
		return nil
	}
	x := -1
	for i, e := range q.paths {
		if x == -1 {
			if !e.executed {
				x = i
			}
		} else {
			if !e.executed && e.rel.Weight() > q.paths[x].rel.Weight() {
				x = i
			}
		}
	}
	if x == -1 {
		return nil
	}
	return q.paths[x]
}

func (q Queue) Get(ix int) (*Path, error) {
	if ix <0 || ix >= len(q.paths) {
		return nil, errors.New("Out of index: " + fmt.Sprint(ix))
	}
	return q.paths[ix], nil
}

func (q *Queue) AddClone(s *Path) (*Path) {
	n := Path{
		rel: s.rel,
		executed: false,
		outcome: kkk.CE_OUTCOME_NULL,
		changed: false,
		cycle: -1,
	}
	for _, ss := range q.paths {
		if !ss.executed && ss.rel.Equals(n.rel) {
			// such type of state is already in queue ready to be executed
			return nil
		}
	}
	q.paths = append(q.paths, &n)
	return &n
}

func (q Queue) Print() string {
	output := ""
	for _, s := range q.paths {
		output += "Cycle        : " + fmt.Sprint(s.cycle) + "\n"
		output += "Cause-Effect : " + s.rel.Print() + "\n"
		output += "Outcome      : " + fmt.Sprint(s.outcome) + "\n"
		output += "Changed      : " + fmt.Sprint(s.changed) + "\n"
		output += "Inputs\n" + s.input.Print() + "\n"
		output += "Outputs\n" + s.output.Print() + "\n"
		output += "\n"
	}
	return output
}

// This method checks if at least one executed state reached (by this ouput) the context of
// the given state.
// This method is used by the reasoner to avoid loop into the graph.
func (q Queue) CheckContext(e *Path) bool {
	for _, s := range q.paths {
		if s.executed && !s.rel.Equals(e.rel) && e.output.PartOf(s.output) {
			return true
		}
	}
	return false
}

func (q *Queue) SortByCycle() {
	sort.Slice(q.paths, func(i, j int) bool {
		return q.paths[i].cycle < q.paths[j].cycle
	  })
}

func (q *Queue) GetCycles() int {
	cycles := -1
	for _, s := range q.paths {
		if s.cycle > cycles {
			cycles = s.cycle
		}
	}
	return cycles
}