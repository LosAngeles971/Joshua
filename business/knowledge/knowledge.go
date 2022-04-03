/*
Knowledge is a set of events and their relationships in terms of cause-effect bindings.
*/
package knowledge

import (
	"fmt"
	"errors"
	log "github.com/sirupsen/logrus"
)

type Knowledge struct {
	max_cycles int
	events []*Event
}

type KnowledgeOption func(*Knowledge) error

func WithSource(source string) KnowledgeOption {
	return func(k *Knowledge) error {
		ee, err := Compile(source)
		if err != nil {
			return err
		}
		k.events = ee
		return nil
	}
}

func WithEvents(ee []*Event) KnowledgeOption {
	return func(k *Knowledge) error {
		if len(ee) == 0 {
			return fmt.Errorf("knowledge cannot be created empty")
		}
		k.events = ee
		return nil
	}
}

func WithMaxCycles(maxCycles int) KnowledgeOption {
	return func(k *Knowledge) error {
		k.max_cycles = maxCycles
		return nil
	}
}

func NewKnowledge(opts ...KnowledgeOption) (Knowledge, error) {
	k := Knowledge{
		max_cycles: 500,
	}
	for _, opt := range opts {
		err := opt(&k)
		if err != nil {
			return k, err
		}
	}
	return k, nil
}

// getEvent returns the event with the given id if it exists into the knowledge
func (u Knowledge) GetEvent(id string) (*Event, bool) {
	for _, e := range u.events {
		if e.getID() == id {
			return e, true
		}
	}
	return nil, false
}

// return the list of events that are cause of the given event
func (u Knowledge) WhoCause(targetEvent Event) []*Event {
	result := []*Event{}
	for _, event := range u.events {
		if event.CanYouCauseThis(targetEvent) {
			result = append(result, event)
		}
	}
	return result
}

//This method returns all possibile paths (given a Knowledge) that end to the given event.
func (k Knowledge) GetAllPathsToEvent(effect *Event) []*Path {
	s := &Stack{}
	discovered := []*Path{}
	for _, cause := range k.WhoCause(*effect) {
		p := createPath(cause, effect)
		discovered = append(discovered, &p)
		s.push(&p)
	}
	for s.size() > 0 {
		p, ok := s.pop()
		if !ok {
			return discovered
		}
		backward(p, k, s)
	}
	return discovered
}

// This function creates a Queue containing all possibile paths to the given event
func (k Knowledge) CreateQueue(data State, effect *Event) Queue {
	return Queue{
		Paths: k.GetAllPathsToEvent(effect),
	}
}

// IsItGoingToHappen verifies if an event occurs given a knowledge and an initial state
func (k Knowledge) isItGoingToHappen(state State, effect *Event) (string, Queue, error) {
	log.Debugf("starting graph analysis with max_cycles to %v", k.max_cycles)
	queue := k.CreateQueue(state, effect)
	log.Tracef("starting queue got the size of %v", queue.Size())
	cycles := 0
	current := state.Clone()
	for {
		path := queue.choose()
		if path == nil {
			log.Debugf("chosen path from queue is nil -> exit with: %v", EFFECT_OUTCOME_EFFECT_FALSE)
			return EFFECT_OUTCOME_EFFECT_FALSE, queue, nil
		}
		err := path.run(*current, cycles)
		if err != nil {
			return EFFECT_OUTCOME_ERROR, queue, err
		}
		log.Tracef("path's outcome -> %v", path.Outcome)
		switch path.Outcome {
		case EFFECT_OUTCOME_LOOP:
			return EFFECT_OUTCOME_ERROR, queue, errors.New("executed path got a loop condition")
		case EFFECT_OUTCOME_ERROR:
			return EFFECT_OUTCOME_ERROR, queue, errors.New("executed path got an error")
		case EFFECT_OUTCOME_UNKNOWN:
			path.Executed = true
			path.Outcome = EFFECT_OUTCOME_UNKNOWN
		case EFFECT_OUTCOME_NULL:
			return EFFECT_OUTCOME_ERROR, queue, errors.New("executed path got no result")
		case EFFECT_OUTCOME_TRUE:
			return EFFECT_OUTCOME_TRUE, queue, nil // the effect happened
		case EFFECT_OUTCOME_EFFECT_FALSE:
			// if the state dit not change the context it does not make sense to have it again into the queue,
			// because it will never reach the desired effect neither it will change the context.
			if path.Changed {
				var pp Path = *path
				if !queue.checkRecurrentOutput(path) {
					// Since the state changed the context by its cause, previous already executed states which are influenced by the
					// execution of this state must be cloned into the queue. This action makes sense only if the current state did not reached
					// a context already reached by another state into the past.
					for _, ppp := range queue.Paths {
						if ppp.Executed {
							ok, err := ppp.isInfluencedBy(pp)
							if err != nil {
								return EFFECT_OUTCOME_ERROR, queue, err
							}
							if ok {
								queue.addPath(ppp)
							}
						}
					}
					// Since the state changed the context by its cause, it does make sense to have an active
					// clone of it into the queue, and update the current globate context.
					// OPEN PROBLEM: should this action be execute outside of the if condition?
					current = path.Output.Clone()
					queue.addPath(path)
				} else {
					// the state reached an already reached context into the past
					// to avoid loopback, the state is not cloned into the queue
					// and it is marked as loop
					path.Outcome = EFFECT_OUTCOME_LOOP
				}
			}
		}
		cycles++
		if cycles > k.max_cycles {
			return EFFECT_OUTCOME_ERROR, queue, fmt.Errorf("reached max cycles -> %v", k.max_cycles)
		}
	}
}