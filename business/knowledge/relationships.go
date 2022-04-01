package knowledge

// A relationship represents the cause-effect binding between two events
type Relationship struct {
	Name   string  // name of the effect's event
	Weight float64 // if the cause's event occurr there is a "weight" probability that the effect occurs
	Effect *Event  // effect's event
}

type RelationshipOption func(*Relationship)

func WithWeight(w float64) RelationshipOption {
	return func(r *Relationship) {
		r.Weight = w
	}
}

func NewRelationship(name string, opts ...RelationshipOption) *Relationship {
	r := &Relationship{
		Name: name,
		Weight: 1.0,
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func (r Relationship) GetWeight() float64 {
	if r.Weight > 1.0 {
		return 1.0
	}
	if r.Weight < 0.0 {
		return 0.0
	}
	return r.Weight
}
