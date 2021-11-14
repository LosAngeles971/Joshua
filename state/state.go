package state

import (
	"fmt"
	"sort"
)

type State struct {
	Data map[string]*Variable `data:"id"`
}

func New() *State {
	return &State{
		Data: map[string]*Variable{},
	}
}

func (c *State) Declare(vname string, vtype string) error {
	var vv Variable
	switch vtype {
	case TYPE_BOOL:
		vv = CreateBool(vname)
	case TYPE_NUMBER:
		vv = CreateNumber(vname)
	default:
		return fmt.Errorf("unrecognized type %v for %v", vtype, vname)
	}
	c.Data[vv.GetName()] = &vv
	return nil
}

func (c *State) Add(vname string, value interface{}) error {
	var vv Variable
	switch value.(type) {
	case bool:
		vv = CreateBool(vname)
		vv.SetValue(value)
	case float64:
		vv = CreateNumber(vname)
		vv.SetValue(value)
	default:
		return fmt.Errorf("unrecognized type of %v as value for %v", value, vname)
	}
	c.Data[vv.GetName()] = &vv
	return nil
}

func (c *State) Update(vname string, v interface{}) error {
	vv, ok := c.Get(vname)
	if ok {
		return vv.SetValue(v)
	} else {
		return c.Add(vname, v)
	}
}

func (c State) Vars() []string {
	keys := []string{}
	for k := range c.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (c State) Contains(k string) bool {
	_, ok := c.Data[k]
	return ok
}

func (c State) Get(k string) (*Variable, bool) {
	v, ok :=  c.Data[k]
	return v, ok
}

// check if child is a subset of father
func (child State) IsSubsetOf(father State) bool {
	for vname, vv := range child.Data {
		v2, ok := father.Get(vname)
		if !ok {
			return false
		}
		if !vv.Equals(*v2) {
			return false
		}
	}
	return true
}

func (c *State) Size() int {
	return len(c.Data)
}

func (c State) Clone() (State) {
	clone := State{
		Data: map[string]*Variable{},
	}
	for vname, v := range c.Data {
		clone.Data[vname] = v.Clone()
	}
	return clone
}

func (c State) Translate() map[string]interface{} {
	vars := map[string]interface{}{}
	for _, v := range c.Data {
		value, ok := v.GetValue()
		if ok { 
			vars[v.GetName()] = value
		}
	}
	return vars
}

/*
func (c State) Print() string {
	output := ""
	keys := make([]string, 0, len(c.Data))
	for k := range c.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		output += k + " = " + fmt.Sprint(c.Data[k].Value) + " defined: " + fmt.Sprint(c.Data[k].Defined) + "\n"
	}
	return output
}
*/