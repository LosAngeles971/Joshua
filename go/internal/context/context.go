package context

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"strconv"
	"math/rand"
)

type Set struct {
	Def		string
	Values 	[]float64
}

func (g Set) Get() float64 {
    return g.Values[rand.Intn(len(g.Values))]
}

func (g Set) Size() int {
    return len(g.Values)
}

type Range interface {
	Size() int
	Get() float64
}

func ParseRange(r string) (Range, error) {
	if strings.HasPrefix(r, "{") && strings.HasSuffix(r, "}") {
		rr := Set{
			Def: r,
			Values: []float64{},
		}
		for _, v := range strings.Split(r[1:len(r)-1], ",") {
			f, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
			if err != nil {
				return nil, err
			}
			rr.Values = append(rr.Values, f)
		}
		return rr, nil
	} else {
		return nil, errors.New("Unrecognized range, it misses braces {}")
	}
}

type Variable struct {
	Name 	string
	Value 	float64
	Defined	bool
	Range	Range
}

func (v Variable) Clone() *Variable {
	vv := Variable{}
	vv.Name = v.Name
	vv.Value = v.Value
	vv.Defined = v.Defined
	vv.Range = v.Range
	return &vv
}

type State struct {
	data map[string]*Variable
}

func Create() State {
	ctx := State{
		data: map[string]*Variable{},
	}
	return ctx
}

func (c State) Clone() (State) {
	clone := Create()
	for k, v := range c.data {
		clone.data[k] = v.Clone()
	}
	return clone
}

func (c State) State() map[string]interface{} {
	state := map[string]interface{}{}
	for _, v := range c.data {
		if v.Defined { 
			state[v.Name] = v.Value
		}
	}
	return state
}

func (c State) PartOf(s State) bool {
	for k, v := range c.data {
		if v.Defined {
			vv, ok := s.data[k]
			if !ok {
				return false
			}
			// why must it be even equals in terms of values????
			// TO AVOID LOOP DUE TO UNCHANGED DATA
			if v.Value != vv.Value {
				return false
			}
		}
	}
	return true
}

func (c *State) Add(v *Variable) {
	c.data[v.Name] = v
}

func (c *State) Update(k string, v float64) {
	vv, ok := c.data[k]
	if ok {
		vv.Defined = true
		vv.Value = v
	} else {
		vv := Variable{
			Name: k,
			Value: v,
			Defined: true,
		}
		c.Add(&vv)
	}
}

func (c *State) Size() int {
	return len(c.data)
}

func (c State) Contains(k string) bool {
	_, ok := c.data[k]
	return ok
}

func (c State) Get(k string) (*Variable, bool) {
	v, ok :=  c.data[k]
	return v, ok
}

func (c State) Print() string {
	output := ""
	keys := make([]string, 0, len(c.data))
	for k := range c.data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		output += k + " = " + fmt.Sprint(c.data[k].Value) + "\n"
	}
	return output
}