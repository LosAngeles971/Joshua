/*
An Event is something that can occur.
An event includes:

- the name
- a list of premises in terms of temporary assignements
- a list of consequences in terms of assignements
- a list of consequences in terms of other events

From a cause/effect perspective, a Event may be a casue and/or an effect.

From model perspective a Event is the vertex of directed graph, where the edges mean cause/effect relationships.
*/
package knowledge

import (
	"fmt"

	"github.com/Knetic/govaluate"
	log "github.com/sirupsen/logrus"
)

const (
	EVENT_OUTCOME_TRUE    = "true"
	EVENT_OUTCOME_FALSE   = "false"
	EVENT_OUTCOME_UNKNOWN = "missing data"
	EVENT_OUTCOME_ERROR   = "error"

	EFFECT_OUTCOME_TRUE         = "true"
	EFFECT_OUTCOME_CAUSE_FALSE  = "cause not happened"
	EFFECT_OUTCOME_EFFECT_FALSE = "effect not happened"
	EFFECT_OUTCOME_UNKNOWN      = "missing data"
	EFFECT_OUTCOME_ERROR        = "error"
	EFFECT_OUTCOME_NULL         = "not verified yet"
	EFFECT_OUTCOME_LOOP         = "true but loop"
)

// Event includes all attribute of an event
type Event struct {
	ID          string                           `yaml:"id"`
	conditions  []*govaluate.EvaluableExpression // list of equivalences used to evaluate if an event can occur
	assignments []Assignment                     // list of assignments executed if an event occurs
	effects     []*Relationship                  // list of effects if the event occurs
}

type EventOption func(*Event)

func WithConditions(c []*govaluate.EvaluableExpression) EventOption {
	return func(e *Event) {
		e.conditions = c
	}
}

func WithAssignments(a []Assignment) EventOption {
	return func(e *Event) {
		e.assignments = a
	}
}

func WithRelationships(r []*Relationship) EventOption {
	return func(e *Event) {
		e.effects = r
	}
}

func NewEvent(id string, opts ...EventOption) *Event {
	e := &Event{
		ID:          id,
		conditions:  []*govaluate.EvaluableExpression{},
		assignments: []Assignment{},
		effects:     []*Relationship{},
	}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

// solveEffects resolves the event's effects starting from the names of the effects
func (event *Event) solveEffects(kkk []*Event) error {
	for _, effect := range event.effects {
		ok := false
		for _, target := range kkk {
			if target.ID == effect.Name {
				effect.Effect = target
				ok = true
			}
		}
		if !ok {
			return fmt.Errorf("effect %v of event %v does not exist", effect.Name, event.ID)
		}
	}
	return nil
}

func (event *Event) addEffects(ee []*Relationship) {
	event.effects = ee
}

func (f Event) getID() string {
	return f.ID
}

// GetWeightTo returns the weight of a specific cause-effect, between the cause and the given effect (if the link exists)
func (e Event) GetWeightTo(effect *Event) float64 {
	w := float64(0.0)
	for _, ef := range e.effects {
		if ef.Effect.getID() == effect.getID() {
			w += ef.Weight
		}
	}
	return w
}

func find(a []string, i string) bool {
	for _, v := range a {
		if v == i {
			return true
		}
	}
	return false
}

// IsInfluencedBy returns true if the given influencer event changes the value of at least one variable used by event's conditions
func (event Event) isInfluencedBy(influencer *Event) (bool, error) {
	vars := []string{}
	for _, a := range influencer.assignments {
		vars = append(vars, a.variable)
	}
	if len(vars) == 0 {
		return false, nil
	}
	for _, s := range event.conditions {
		for _, v := range s.Vars() {
			if find(vars, v) {
				return true, nil
			}
		}
	}
	return false, nil
}

// IsValid returns false if the event does not have resolved effects
func (event Event) IsValid() error {
	for _, e := range event.effects {
		if e.Effect == nil {
			return fmt.Errorf("event %v has the undefined effect %v", event.ID, e.Name)
		}
	}
	return nil
}

// CanYouCauseThis returns true if the event can cause the given effect event
func (event Event) CanYouCauseThis(effectEvent Event) bool {
	for _, effect := range event.effects {
		if effect.Effect.getID() == effectEvent.getID() && effect.GetWeight() > 0.0 {
			return true
		}
	}
	return false
}

// Run executes event's assignements if the event's conditions are all true
func (f *Event) Run(input State) (string, State, error) {
	output := input.Clone()
	for _, expr := range f.conditions {
		log.Tracef("checking condition [%v] of event [%v]", expr, f.ID)
		if !output.AreDefined(expr.Vars()) {
			return EVENT_OUTCOME_UNKNOWN, *output, nil
		}
		result, err := expr.Evaluate(output.Translate())
		log.Debugf("condition [%v] of event [%v] got result [%v] and error [%v]", expr, f.ID, result, err)
		if err != nil {
			return EVENT_OUTCOME_ERROR, *output, err
		}
		vv, ok := result.(bool)
		if !ok {
			return EVENT_OUTCOME_ERROR, *output, fmt.Errorf("condition must be boolean [%v]", expr.String())
		}
		if !vv {
			return EVENT_OUTCOME_FALSE, *output, nil
		}
	}
	for _, expr := range f.assignments {
		log.Tracef("running assigment [%v] of event [%v]", expr, f.ID)
		if !output.AreDefined(expr.expr.Vars()) {
			return EVENT_OUTCOME_UNKNOWN, *output, nil
		}
		result, err := expr.expr.Evaluate(output.Translate())
		log.Debugf("assignment [%v] of event [%v] got result [%v] and error [%v]", expr, f.ID, result, err)
		if err != nil {
			return EVENT_OUTCOME_ERROR, *output, err
		}
		err = output.Update(expr.variable, result)
		if err != nil {
			return EVENT_OUTCOME_ERROR, *output, err
		}
	}
	return EVENT_OUTCOME_TRUE, *output, nil
}
