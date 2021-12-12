package knowledge

import (
	"io/ioutil"
	"sort"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Queue struct {
	Paths []*Path	`yaml:"queue"`
}

// Serialize creates a YAML representation of the queue
// This method is useful to print the final queue, the latter explains the reasoning
func (q Queue) Serialize() (string, error) {
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
	qq.sortByCycle()
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

// Size returns the number of Paths into the Queue
func (q Queue) Size() (int) {
	return len(q.Paths)
}

// AddPath clones the given path in order to add it into the queue
// Such situazion happens, when the state changed and the change makes
// an already executed path applicable again
func (q *Queue) addPath(s *Path) (*Path) {
	log.Trace("cloning path...")
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
func (q Queue) checkRecurrentOutput(e *Path) bool {
	for _, s := range q.Paths {
		if s.Executed && !s.equals(e) && e.Output.IsSubsetOf(s.Output) {
			return true
		}
	}
	return false
}

func (q *Queue) sortByCycle() {
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

// Choose returns the current, best Path to the success event
func (q Queue) choose() (*Path) {
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
			if !p.Executed && p.getWeight() > q.Paths[x].getWeight() {
				x = i
			}
		}
	}
	if x == -1 {
		return nil
	}
	return q.Paths[x]
}