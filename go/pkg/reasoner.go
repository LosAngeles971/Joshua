package pkg

import (
	"errors"
	kkk "it/losangeles971/joshua/internal/knowledge"
	ctx "it/losangeles971/joshua/internal/context"
	"it/losangeles971/joshua/internal/runtime"
	"text/tabwriter"
	"fmt"
	"os"
)

/*
Given a hypothetical effect, is it happening? depending on the current Universe of knowledge+context?
If after a reasoning at least one cause-effect relationship changed its state, you need to do a new cycle
of reasoning
*/
func Verify(k kkk.Knowledge, init ctx.State, effect kkk.Event, max_cycles int) (string, runtime.Queue, error) {
	queue := runtime.Queue{}
	queue.Populate(init, k, effect)
	cycles := 0
	current := init.Clone()
	for true {
		state := queue.Choose()
		if state == nil {
			return kkk.CE_OUTCOME_EFFECT_FALSE, queue, nil
		}
		err := state.Run(current, cycles)
		if err != nil {
			return kkk.CE_OUTCOME_ERROR, queue, err
		}
		switch state.Outcome() {
		case kkk.CE_OUTCOME_LOOP:
			return kkk.CE_OUTCOME_ERROR, queue, errors.New("Executed state in loop condition")
		case kkk.CE_OUTCOME_ERROR:
			return kkk.CE_OUTCOME_ERROR, queue, errors.New("Executed state in error condition")
		case kkk.CE_OUTCOME_UNKNOWN:
			return kkk.CE_OUTCOME_ERROR, queue, errors.New("Executed state in unknown condition")
		case kkk.CE_OUTCOME_NULL:
			return kkk.CE_OUTCOME_ERROR, queue, errors.New("Executed state in nil condition")
		case kkk.CE_OUTCOME_TRUE:
			return kkk.CE_OUTCOME_TRUE, queue, nil // the effect happened
		case kkk.CE_OUTCOME_EFFECT_FALSE:
			// if thes state dit not change the context it does not make sense to have it again into the queue,
			// because it will never reach the desired effect neither it will change the context.
			if state.Changed() {
				if !queue.CheckContext(state) {
					// Since the state changed the context by its cause, previous already executed states which are influenced by the 
					// execution of this state must be cloned into the queue. This action makes sense only if the current state did not reached
					// a context already reached by another state into the past.
					for i := 0; i < queue.Size(); i++ {
						s, _ := queue.Get(i)
						if s.Executed() {
							ok, err := s.IsInfluenced(state.Relationship())
							if err != nil {
								return kkk.CE_OUTCOME_ERROR, queue, err
							}
							if ok {
								queue.AddClone(s)
							}
						}
					}
					// Since the state changed the context by its cause, it does make sense to have an active
					// clone of it into the queue, and update the currenct globate context.
					// OPEN PROBLEM: should this action be execute outside of the if condition?
					current = state.Output().Clone()
					queue.AddClone(state)
				} else {
					// the state reached an already reached context into the past
					// to avoid loopback, the state is not cloned into the queue
					// and it is marked as loop
					state.Loop()
				}
			}
		}
		cycles++
		if cycles > max_cycles {
			return kkk.CE_OUTCOME_ERROR, queue, errors.New("Reached max cycles")
		}
	}
	return kkk.CE_OUTCOME_ERROR, queue, errors.New("Cycle broken")
}

func PrintSummary(outcome string, queue runtime.Queue) {
	const padding = 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintln(w, "Outcome\t" + outcome + "\t")
	fmt.Fprintln(w, "Cycles\t" + fmt.Sprint(queue.GetCycles()) + "\t")
	fmt.Fprintln(w, "Queue's size\t" + fmt.Sprint(queue.Size()) + "\t")
	w.Flush()
}

func PrintQueue(queue runtime.Queue) {
	const padding = 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', tabwriter.AlignRight|tabwriter.Debug)
	queue.SortByCycle()
	for i := 0; i < queue.Size(); i++ {
		s, _ := queue.Get(i) 
		fmt.Fprintln(w, "Cycle\t" + fmt.Sprint(s.Cycle()) + "\t")
		fmt.Fprintln(w, "Cause\t" + fmt.Sprint(s.Relationship().Cause.GetID()) + "\t")
		fmt.Fprintln(w, "Effect\t" + fmt.Sprint(s.Relationship().Effect.GetID()) + "\t")
		fmt.Fprintln(w, "Outcome\t" + fmt.Sprint(s.Outcome()) + "\t")
		fmt.Fprintln(w, "Context changed\t" + fmt.Sprint(s.Changed()) + "\t")
		fmt.Fprintln(w, "\t\t")
	}
	w.Flush()
}