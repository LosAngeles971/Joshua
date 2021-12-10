/*
This package exports the main functions of Joshua to the CLI and/or external programs.

The main function is "IsItGoingToHappen"::
Given an event, named success event, the function verifies it may happen starting from an initial state.
*/
package engine

import (
	"errors"
	"fmt"
	"it/losangeles971/joshua/knowledge"
	"it/losangeles971/joshua/state"
	"os"
	"text/tabwriter"
)

// PrintSummary prints the results of an execution to the console (standard output)
func PrintSummary(outcome string, queue knowledge.Queue) {
	const padding = 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintln(w, "Outcome\t" + outcome + "\t")
	fmt.Fprintln(w, "Cycles\t" + fmt.Sprint(queue.GetCycles()) + "\t")
	fmt.Fprintln(w, "Queue's size\t" + fmt.Sprint(queue.Size()) + "\t")
	w.Flush()
}

// isItGoingToHappen verifies if an event occurs given a knowledge and an initial state
func isItGoingToHappen(k knowledge.Knowledge, init state.State, effect *knowledge.Event, max_cycles int) (string, knowledge.Queue, error) {
	queue := k.CreateQueue(init, effect)
	cycles := 0
	current := init.Clone()
	for {
		path := queue.Choose()
		if path == nil {
			return knowledge.EFFECT_OUTCOME_EFFECT_FALSE, queue, nil
		}
		err := path.Run(current, cycles)
		if err != nil {
			return knowledge.EFFECT_OUTCOME_ERROR, queue, err
		}
		switch path.Outcome {
		case knowledge.EFFECT_OUTCOME_LOOP:
			return knowledge.EFFECT_OUTCOME_ERROR, queue, errors.New("executed path got a loop condition")
		case knowledge.EFFECT_OUTCOME_ERROR:
			return knowledge.EFFECT_OUTCOME_ERROR, queue, errors.New("executed path got an error")
		case knowledge.EFFECT_OUTCOME_UNKNOWN:
			path.Executed = true
			path.Outcome = knowledge.EFFECT_OUTCOME_UNKNOWN 
		case knowledge.EFFECT_OUTCOME_NULL:
			return knowledge.EFFECT_OUTCOME_ERROR, queue, errors.New("executed path got no result")
		case knowledge.EFFECT_OUTCOME_TRUE:
			return knowledge.EFFECT_OUTCOME_TRUE, queue, nil // the effect happened
		case knowledge.EFFECT_OUTCOME_EFFECT_FALSE:
			// if the state dit not change the context it does not make sense to have it again into the queue,
			// because it will never reach the desired effect neither it will change the context.
			if path.Changed {
				var pp knowledge.Path = *path
				if !queue.CheckRecurrentOutput(path) {
					// Since the state changed the context by its cause, previous already executed states which are influenced by the 
					// execution of this state must be cloned into the queue. This action makes sense only if the current state did not reached
					// a context already reached by another state into the past.
					for _, ppp := range queue.Paths {
						if ppp.Executed {
							ok, err := ppp.IsInfluencedBy(pp)
							if err != nil {
								return knowledge.EFFECT_OUTCOME_ERROR, queue, err
							}
							if ok {
								queue.AddPath(ppp)
							}
						}
					}
					// Since the state changed the context by its cause, it does make sense to have an active
					// clone of it into the queue, and update the current globate context.
					// OPEN PROBLEM: should this action be execute outside of the if condition?
					current = path.Output.Clone()
					queue.AddPath(path)
				} else {
					// the state reached an already reached context into the past
					// to avoid loopback, the state is not cloned into the queue
					// and it is marked as loop
					path.Outcome = knowledge.EFFECT_OUTCOME_LOOP
				}
			}
		}
		cycles++
		if cycles > max_cycles {
			return knowledge.EFFECT_OUTCOME_ERROR, queue, errors.New("reached max cycles")
		}
	}
}

// IsItGoingToHappen exports the internal function isItGoingToHappen
func IsItGoingToHappen(ksource string, init state.State, effect string, max_cycles int) (string, knowledge.Queue, error) {
	kkk, err := knowledge.Load(ksource)
	if err != nil {
		return "", knowledge.Queue{}, err
	}
	ee, ok := kkk.GetEvent(effect)
	if !ok {
		return "", knowledge.Queue{}, fmt.Errorf("success event %v does not exist into the knowledge", effect)
	}
	return isItGoingToHappen(kkk, init, ee, max_cycles)
}