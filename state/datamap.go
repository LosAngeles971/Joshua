package state

import (
	"fmt"
	"sort"
)

type Datamap struct {
	Vars map[string]interface{}
}

type SimpleState struct {
	Data *Datamap
}

func NewSimpleState() SimpleState {
	d := Datamap{
		Vars: map[string]interface{}{},
	}
	return SimpleState{
		Data: &d,
	}
}

func (s SimpleState) Add(name string, value interface{}) error {
	switch value.(type) {
	case bool, float64:
		s.Data.Vars[name] = value
		return nil
	default:
		return fmt.Errorf("type of data not handled: %t", value)
	}
}

func (s SimpleState) Update(name string, value interface{}) error {
	return s.Add(name, value)
}

func (s SimpleState) Vars() []string {
	keys := []string{}
	for k := range s.Data.Vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (s SimpleState) IsDefined(name string) bool {
	_, ok := s.Data.Vars[name]
	return ok
}

func (s SimpleState) Get(name string) (interface{}, bool) {
	v, ok := s.Data.Vars[name]
	return v, ok
}

func (s SimpleState) GetType(name string) (int, error) {
	if !s.IsDefined(name) {
		return 0, fmt.Errorf("variable %s is not defined", name)
	}
	v, _ := s.Get(name)
	switch v.(type) {
	case bool:
		return TYPE_BOOL, nil
	case float64:
		return TYPE_NUMBER, nil
	default:
		return 0, fmt.Errorf("variable %s is of not supported type %t", name, v)
	}
}

// check if child is a "pure" subset of father
func (child SimpleState) IsSubsetOf(father State) bool {
	for name, value := range child.Data.Vars {
		vv, ok := father.Get(name)
		if !ok {
			return false
		}
		if vv != value {
			return false
		}
	}
	return true
}

func (s SimpleState) Size() int {
	return len(s.Data.Vars)
}

func (s SimpleState) Clone() State {
	clone := NewSimpleState()
	for name, value := range s.Data.Vars {
		clone.Add(name, value)
	}
	return clone
}

func (s SimpleState) Translate() map[string]interface{} {
	return s.Data.Vars
}