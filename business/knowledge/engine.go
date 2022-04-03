/*
This package exports the main functions of Joshua to the CLI and/or external programs.

The main function is "IsItGoingToHappen"::
Given an event, named success event, the function verifies it may happen starting from an initial state.
*/
package knowledge

import (
	"fmt"
	"os"
	"text/tabwriter"

	log "github.com/sirupsen/logrus"
)

type Solution struct {
	Outcome string
	Chain   Queue
	Err     error
}

type Engine struct {
	kkk       Knowledge
}

func NewEngine(source string, maxCycles int) (Engine, error) {
	kkk, err := NewKnowledge(WithSource(source), WithMaxCycles(maxCycles))
	if err != nil {
		return Engine{}, err
	}
	return Engine{
		kkk: kkk,
	}, nil
}

func NewSolution(o string, q Queue, e error) Solution {
	return Solution{
		Outcome: o,
		Chain:   q,
		Err:     e,
	}
}

// PrintSummary prints the results of an execution to the console (standard output)
func (s Solution) PrintSummary() {
	const padding = 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintln(w, "Outcome\t"+s.Outcome+"\t")
	fmt.Fprintln(w, "Cycles\t"+fmt.Sprint(s.Chain.GetCycles())+"\t")
	fmt.Fprintln(w, "Queue's size\t"+fmt.Sprint(s.Chain.Size())+"\t")
	w.Flush()
}

func (s Solution) PrintFullChain() {
	y, err := s.Chain.Serialize()
	if err != nil {
		log.Errorf("error while serializing chain -> %v", s.Err)
	} else {
		fmt.Println(y)
	}
}

func (s Solution) PrintChain() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintln(w, "Cycle\tCause\tEffect\tOutcome\t")
	for _, path := range s.Chain.Paths {
		if path.Executed {
			for _, edge := range path.Path {
				if edge.Outcome != EFFECT_OUTCOME_CAUSE_FALSE {
					fmt.Fprintf(w, "%v\t%s\t%s\t%s\t\n", path.Cycle, edge.Cause.ID, edge.Effect.ID, edge.Outcome)
				}
			}
		}
	}
	w.Flush()
}

// IsItGoingToHappen exports the internal function isItGoingToHappen
func (engine Engine) IsItGoingToHappen(state State, effect string) Solution {
	ee, ok := engine.kkk.GetEvent(effect)
	if !ok {
		return NewSolution("", Queue{}, fmt.Errorf("success event %v does not exist into the knowledge", effect))
	}
	return NewSolution(engine.kkk.isItGoingToHappen(state, ee))
}
