package knowledge

import (
	"errors"
	"sort"
	"io/ioutil"
	ctx "it/losangeles971/joshua/internal/context"
	"gopkg.in/yaml.v2"
)

func (s Path) isInfluencedBy(p Path) (bool, error) {
	for _, p_e := range p.Path {
		for _, s_e := range s.Path {
			if ok, err := s_e.IsInfluencedBy(p_e); ok || err != nil {
				return ok, err
			}
		}
	}
	return false, nil
}

// This function run a path starting from a specific index
func (p *Path) run(input ctx.State, tail int) (error) {
	for i := tail; i >= 0; i-- {
		e := p.Path[i]
		outcome, output, err := e.Cause.EffectHappen(input, e.Effect)
		if err != nil {
			return err
		}
		p.outcome = outcome
		p.output = output.Clone()
		if !p.input.PartOf(p.output) {
			p.changed = true
		}
	}
	return nil
}

/* 
Argument input comes from the execution of previous paths
and override the internal input of the receiver's struct.
The verification must got for attempts:
- just the first rel
- if not the first -1 and then the first
- ...
- if not the first - n, first - n -1, ..., the the first
*/
func (p *Path) Run(input ctx.State, cycle int) (error) {
	p.cycle = cycle
	if p.executed {
		return errors.New("Asked to run an already executed state")
	}
	p.input = input.Clone()
	p.executed = true
	p.changed = false
	for i := 0; i < len(p.Path); i++ {
		err := p.run(input, i)
		if err != nil {
			return err
		}
		input = p.output.Clone()
		switch p.outcome {
		case CE_OUTCOME_CAUSE_FALSE:
			return nil;
		case CE_OUTCOME_EFFECT_FALSE:
			return nil
		case CE_OUTCOME_ERROR:
			return errors.New("outcome is error")
		case CE_OUTCOME_LOOP:
		case CE_OUTCOME_NULL:
			return errors.New("outcome is null")
		case CE_OUTCOME_TRUE:
		case CE_OUTCOME_UNKNOWN:
			return nil
		}
	}
	return nil
}

type Queue struct {
	paths []*Path	`yaml:"queue"`
}

/* 
This method serialize a queue of paths into a yaml.
At the end of the reasoning process, the queue represents the found solution,
thus this method can be used to save the solution as yaml string.
*/
func (q Queue) Solution() (string, error) {
	s, err := yaml.Marshal(q)
	if err != nil {
		return "", err
	}
	return string(s), nil
}

/* 
This method serialize a queue of paths into a yaml.
At the end of the reasoning process, the queue represents the found solution,
thus this method can be used to save the solution as yaml file.
*/
func (q Queue) Serialize(filename string) (error) {
	d, err := yaml.Marshal(q)
	if err != nil {
		return err
	}
	ioutil.WriteFile(filename, d, 0644)
	return nil
}

func (q Queue) Size() (int) {
	return len(q.paths)
}

/*
func (q Queue) FindByRelationship(rel kkk.Relationship) (*Path) {
	for _, p := range q.paths {
		for i := range p.rel {
			if p.rel[i].Equals(rel) {
				return p
			}
		}
	}
	return nil
}
*/

func (q *Queue) addClone(s *Path) (*Path) {
	n := Path{
		Path: s.Path,
		executed: false,
		outcome: CE_OUTCOME_NULL,
		changed: false,
		cycle: -1,
	}
	for _, ss := range q.paths {
		if !ss.executed && ss.Equals(&n) {
			// such type of path is already in queue ready to be executed
			return nil
		}
	}
	q.paths = append(q.paths, &n)
	return &n
}

/*
This method checks if the given executed path produced an output 
already reached by another previously executed path into the queue.
In case the output of given path has been already reached by others,
the given path can be considered a potential loop and it can be 
discarderd even if it changed the state.
*/
func (q Queue) CheckContext(e *Path) bool {
	for _, s := range q.paths {
		if s.executed && !s.Equals(e) && e.output.PartOf(s.output) {
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

// This function returns the number of cycles executed by the engine
func (q *Queue) GetCycles() int {
	cycles := -1
	for _, s := range q.paths {
		if s.cycle > cycles {
			cycles = s.cycle
		}
	}
	return cycles
}

// This function choose the Path to check for solving the problem
func (q Queue) Choose() (*Path) {
	if len(q.paths) < 1 {
		return nil
	}
	x := -1
	for i, p := range q.paths {
		if x == -1 {
			if !p.executed {
				x = i
			}
		} else {
			if !p.executed && p.GetWeight() > q.paths[x].GetWeight() {
				x = i
			}
		}
	}
	if x == -1 {
		return nil
	}
	return q.paths[x]
}

// This function creates a Queue to solve a problem
func createQueue(data ctx.State, k Knowledge, effect *Event) Queue {
	q := Queue{
		paths: GetAllPaths(k, effect),
	}
	return q
}

/*
Given a hypothetical effect, is it happening? depending on the current Universe of knowledge+context?
If after a reasoning at least one cause-effect relationship changed its state, you need to do a new cycle
of reasoning
*/
func Reason(k Knowledge, init ctx.State, effect *Event, max_cycles int) (string, Queue, error) {
	queue := createQueue(init, k, effect)
	cycles := 0
	current := init.Clone()
	for true {
		path := queue.Choose()
		if path == nil {
			return CE_OUTCOME_EFFECT_FALSE, queue, nil
		}
		err := path.Run(current, cycles)
		if err != nil {
			return CE_OUTCOME_ERROR, queue, err
		}
		switch path.outcome {
		case CE_OUTCOME_LOOP:
			return CE_OUTCOME_ERROR, queue, errors.New("Executed state in loop condition")
		case CE_OUTCOME_ERROR:
			return CE_OUTCOME_ERROR, queue, errors.New("Executed state in error condition")
		case CE_OUTCOME_UNKNOWN:
			return CE_OUTCOME_ERROR, queue, errors.New("Executed state in unknown condition")
		case CE_OUTCOME_NULL:
			return CE_OUTCOME_ERROR, queue, errors.New("Executed state in nil condition")
		case CE_OUTCOME_TRUE:
			return CE_OUTCOME_TRUE, queue, nil // the effect happened
		case CE_OUTCOME_EFFECT_FALSE:
			// if thes state dit not change the context it does not make sense to have it again into the queue,
			// because it will never reach the desired effect neither it will change the context.
			if path.changed {
				var pp Path = *path
				if !queue.CheckContext(path) {
					// Since the state changed the context by its cause, previous already executed states which are influenced by the 
					// execution of this state must be cloned into the queue. This action makes sense only if the current state did not reached
					// a context already reached by another state into the past.
					for _, ppp := range queue.paths {
						if ppp.executed {
							ok, err := ppp.isInfluencedBy(pp)
							if err != nil {
								return CE_OUTCOME_ERROR, queue, err
							}
							if ok {
								queue.addClone(ppp)
							}
						}
					}
					// Since the state changed the context by its cause, it does make sense to have an active
					// clone of it into the queue, and update the currenct globate context.
					// OPEN PROBLEM: should this action be execute outside of the if condition?
					current = path.output.Clone()
					queue.addClone(path)
				} else {
					// the state reached an already reached context into the past
					// to avoid loopback, the state is not cloned into the queue
					// and it is marked as loop
					path.outcome = CE_OUTCOME_LOOP
				}
			}
		}
		cycles++
		if cycles > max_cycles {
			return CE_OUTCOME_ERROR, queue, errors.New("Reached max cycles")
		}
	}
	return CE_OUTCOME_ERROR, queue, errors.New("Cycle broken")
}