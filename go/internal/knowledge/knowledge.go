package knowledge

type Knowledge struct {
	Events 				[]Event
	Relationships		[]Relationship
}

// This function search for a Event
func (u Knowledge) GetEvent(id string) (Event, bool) {
	for _, f := range u.Events {
		if f.GetID() == id {
			return f, true
		}
	}
	return Event{}, false
}

// This function search for a Event
func (u Knowledge) GetRelationship(cause Event, effect Event) (Relationship, bool) {
	for _, r := range u.Relationships {
		if r.Cause.GetID() == cause.GetID() && r.Effect.GetID() == effect.GetID() {
			return r, true
		}
	}
	return Relationship{}, false
}

// This function search for a Event
func (u Knowledge) GetEffects(cause Event) ([]Event) {
	effects := []Event{}
	for _, r := range u.Relationships {
		if r.Cause.GetID() == cause.GetID() {
			effects = append(effects, r.Effect)
		}
	}
	return effects
}

func (u Knowledge) IsCauseOf(cause Event) ([]Relationship) {
	result := []Relationship{}
	for _, r := range u.Relationships {
		if r.Cause.GetID() == cause.GetID() {
			result = append(result, r)
		}
	}
	return result
}

func (u Knowledge) IsEffectOf(effect Event) ([]Relationship) {
	result := []Relationship{}
	for _, r := range u.Relationships {
		if r.Effect.GetID() == effect.GetID() {
			result = append(result, r)
		}
	}
	return result
}