/*
This file provides capabilities to handle states for Joshua programs.
*/
package knowledge

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

const (
	ENC_JSON = 0
	ENC_YAML = 1
)

// A state is a map of values.
type State struct {
	Vars map[string]interface{} `yaml:"vars"`
}

type StateOption func(*State)

func WithData(data []byte, enc int) StateOption {
	return func(s *State) {
		ss := State{}
		var err error
		if enc == ENC_JSON {
			err = json.Unmarshal(data, &ss)
		} else {
			err = yaml.Unmarshal(data, &ss)
		}
		if err != nil {
			return
		}
		for k, v := range ss.Vars {
			s.Vars[k] = v
		}
	}
}

// WithMap loads the state from a given map
func WithMap(mm map[string]interface{}) StateOption {
	return func(s *State) {
		for k, v := range mm {
			s.Vars[k] = v
		}
	}
}

// Create a new state
func NewState(opts ...StateOption) *State {
	s := &State{
		Vars: map[string]interface{}{},
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s State) Add(name string, value interface{}) error {
	switch value.(type) {
	case bool, float64:
		s.Vars[name] = value
		return nil
	default:
		return fmt.Errorf("type of data not handled: %t", value)
	}
}

func (s State) Update(name string, value interface{}) error {
	return s.Add(name, value)
}

func (s State) Get(name string) (interface{}, bool) {
	v, ok := s.Vars[name]
	return v, ok
}

// IsSubsetOf checks if the state is a subset of the given state
func (child State) IsSubsetOf(father State) bool {
	for name, value := range child.Vars {
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

func (s State) Size() int {
	return len(s.Vars)
}

func (s State) Clone() *State {
	clone := NewState(WithMap(s.Vars))
	for name, value := range s.Vars {
		clone.Add(name, value)
	}
	return clone
}

// Translate converts the state into a map
func (s State) Translate() map[string]interface{} {
	return s.Vars
}

func (s State) AreDefined(vv []string) bool {
	for i := range vv {
		_, ok := s.Get(vv[i])
		if !ok {
			return false
		}
	}
	return true
}