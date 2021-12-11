/*
Datamap is a really simple implementation of the state's interface, since it is just a map of interfaces,
where every tuple (k,v) is a variable.
*/
package knowledge

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

// State is a map of current values of all defined variables
type State struct {
	Vars map[string]interface{} `yaml:"vars"`
}

type StateOption func(*State)

func WithYAML(filename string) StateOption {
	return func(s *State) {
		in, err := ioutil.ReadFile(filename)
		if err != nil {
			return
		}
		ss := State{}
		yaml.Unmarshal(in, &ss)
		for k, v := range ss.Vars {
			s.Vars[k] = v
		}
	}
}

func WithJSON(filename string) StateOption {
	return func(s *State) {
		in, err := ioutil.ReadFile(filename)
		if err != nil {
			return
		}
		ss := State{}
		json.Unmarshal(in, &ss)
		for k, v := range ss.Vars {
			s.Vars[k] = v
		}
	}
}

func WithMap(mm map[string]interface{}) StateOption {
	return func(s *State) {
		for k, v := range mm {
			s.Vars[k] = v
		}
	}
}

// Create a new SimpleState with no variables
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

// check if child is a "pure" subset of father
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

func (s State) Translate() map[string]interface{} {
	return s.Vars
}
