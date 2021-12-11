/*
Knowledge is a set of events and their relationships in terms of cause-effect bindings.
*/
package knowledge

type Knowledge struct {
	events []*Event
}

// Load compile a knowledge's source into a knowledge object
func Load(source string) (Knowledge, error) {
	var err error
	kkk := Knowledge{}
	kkk.events, err = parse(source)
	if err != nil {
		return kkk, err
	}
	for _, event := range kkk.events {
		err = event.solveEffects(kkk.events)
		if err != nil {
			return Knowledge{}, err
		}
	}
	return kkk, nil
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