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
)

type Engine struct {
	kkk Knowledge
}

func NewEngine(source string) (Engine, error) {
	kkk, err := NewKnowledge(WithSource(source), WithMaxCycles(500))
	if err != nil {
		return Engine{}, err
	}
	return Engine{
		kkk: kkk,
	}, nil
}

// PrintSummary prints the results of an execution to the console (standard output)
func PrintSummary(outcome string, queue Queue) {
	const padding = 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintln(w, "Outcome\t" + outcome + "\t")
	fmt.Fprintln(w, "Cycles\t" + fmt.Sprint(queue.GetCycles()) + "\t")
	fmt.Fprintln(w, "Queue's size\t" + fmt.Sprint(queue.Size()) + "\t")
	w.Flush()
}

// IsItGoingToHappen exports the internal function isItGoingToHappen
func (engine Engine) IsItGoingToHappen(state State, effect string) (string, Queue, error) {
	ee, ok := engine.kkk.GetEvent(effect)
	if !ok {
		return "", Queue{}, fmt.Errorf("success event %v does not exist into the knowledge", effect)
	}
	return engine.kkk.isItGoingToHappen(state, ee)
}