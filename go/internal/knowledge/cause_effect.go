package knowledge

import (
	ctx "it/losangeles971/joshua/internal/context"
)

const (
	CE_OUTCOME_TRUE         = "true"
	CE_OUTCOME_CAUSE_FALSE  = "cause not happened"
	CE_OUTCOME_EFFECT_FALSE = "effect not happened"
	CE_OUTCOME_UNKNOWN      = "missing data"
	CE_OUTCOME_ERROR        = "error"
	CE_OUTCOME_NULL         = "not verified yet"
	CE_OUTCOME_LOOP         = "true but loop"
)

type Relationship struct {
	Indirect  	string		`yaml:"event"`
	Weight  	float64		`yaml:"weight"`
	Effect 		*Event
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

/*
func (l Relationship) Equals(r Relationship) bool {
	if l.Cause.GetID() == r.Cause.GetID() && l.Effect.GetID() == r.Effect.GetID() {
		return true
	}
	return false
}
*/

func (cause *Event) EffectHappen(init ctx.State, effect *Event) (string, ctx.State, error) {
	data := init.Clone()
	outcome, err := cause.CanHappen(&data)
	if err != nil {
		return CE_OUTCOME_ERROR, data, err
	}
	if outcome == EVENT_OUTCOME_FALSE {
		return CE_OUTCOME_CAUSE_FALSE, data, nil
	}
	if outcome == EVENT_OUTCOME_UNKNOWN {
		return CE_OUTCOME_UNKNOWN, data, nil
	}
	outcome, err = effect.CanHappen(&data)
	if err != nil {
		return CE_OUTCOME_ERROR, data, err
	}
	if outcome == EVENT_OUTCOME_FALSE {
		return CE_OUTCOME_EFFECT_FALSE, data, nil
	}
	if outcome == EVENT_OUTCOME_UNKNOWN {
		return CE_OUTCOME_UNKNOWN, data, nil
	}
	return CE_OUTCOME_TRUE, data, nil
}


/*
func (influenced Relationship) IsInfluencedBy(cause Event, effect Event) (bool, error) {
	ok, err := influenced.Cause.IsInfluencedBy(influencer.Cause)
	if err != nil {
		return ok, err
	}
	if ok {
		return true, nil
	}
	ok, err = influenced.Cause.IsInfluencedBy(influencer.Effect)
	if err != nil {
		return ok, err
	}
	if ok {
		return true, nil
	}
	ok, err = influenced.Effect.IsInfluencedBy(influencer.Cause)
	if err != nil {
		return ok, err
	}
	if ok {
		return true, nil
	}
	ok, err = influenced.Effect.IsInfluencedBy(influencer.Effect)
	if err != nil {
		return ok, err
	}
	if ok {
		return true, nil
	}
	return false, nil
}
*/