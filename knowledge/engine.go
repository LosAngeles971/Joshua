package knowledge

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// IsItGoingToHappen verifies if an event occurs given a knowledge and an initial state
func IsItGoingToHappen(k Knowledge, init State, effect *Event, max_cycles int) (string, Queue, error) {
	log.Debugf("starting graph analysis with max_cycles to %v", max_cycles)
	queue := k.CreateQueue(init, effect)
	log.Tracef("starting queue got the size of %v", queue.Size())
	cycles := 0
	current := init.Clone()
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
		if cycles > max_cycles {
			return EFFECT_OUTCOME_ERROR, queue, fmt.Errorf("reached max cycles -> %v", max_cycles)
		}
	}
}
