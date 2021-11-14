package knowledge

import (
	"io/ioutil"
	"sort"

	"gopkg.in/yaml.v3"
)

type Queue struct {
	Paths []*Path	`yaml:"queue"`
}

/* 
This method serialize a queue of paths into a yaml.
At the end of the reasoning process, the queue represents the found solution,
thus this method can be used to save the solution as yaml string.
*/
func (q Queue) Serialize(clean bool) (string, error) {
	qq := Queue{
		Paths: []*Path{},
	}
	for _, p := range q.Paths {
		if p.Executed {
			if p.Outcome == EFFECT_OUTCOME_EFFECT_FALSE || p.Outcome == EFFECT_OUTCOME_TRUE {
				qq.Paths = append(qq.Paths, p)
			}
		}
	}
	qq.SortByCycle()
	s, err := yaml.Marshal(&qq)
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
func (q Queue) Save(filename string) (error) {
	d, err := yaml.Marshal(&q)
	if err != nil {
		return err
	}
	ioutil.WriteFile(filename, d, 0644)
	return nil
}

func (q Queue) Size() (int) {
	return len(q.Paths)
}

// This method clone an executed Path and put it into the queue
// Such situazion happens, when the state changed and the change makes
// an already executed path applicable again
func (q *Queue) AddPath(s *Path) (*Path) {
	n := s.clone()
	for _, ss := range q.Paths {
		if !ss.Executed && ss.equals(n) {
			// such type of path is already in queue ready to be executed
			return nil
		}
	}
	q.Paths = append(q.Paths, n)
	return n
}

/*
This method checks if the given executed path produced an output 
already reached by another previously executed path into the queue.
In case the output of given path has been already reached by others,
the given path can be considered a potential loop and it can be 
discarderd even if it changed the state.
*/
func (q Queue) CheckRecurrentOutput(e *Path) bool {
	for _, s := range q.Paths {
		if s.Executed && !s.equals(e) && e.Output.IsSubsetOf(s.Output) {
			return true
		}
	}
	return false
}

func (q *Queue) SortByCycle() {
	sort.Slice(q.Paths, func(i, j int) bool {
		return q.Paths[i].Cycle < q.Paths[j].Cycle
	  })
}

// This function returns the number of cycles executed by the engine
func (q *Queue) GetCycles() int {
	cycles := -1
	for _, s := range q.Paths {
		if s.Cycle > cycles {
			cycles = s.Cycle
		}
	}
	return cycles
}

// This function choose the Path to check for solving the problem
func (q Queue) Choose() (*Path) {
	if len(q.Paths) < 1 {
		return nil
	}
	x := -1
	for i, p := range q.Paths {
		if x == -1 {
			if !p.Executed {
				x = i
			}
		} else {
			if !p.Executed && p.GetWeight() > q.Paths[x].GetWeight() {
				x = i
			}
		}
	}
	if x == -1 {
		return nil
	}
	return q.Paths[x]
}