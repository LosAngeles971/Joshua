package state

import (
	"fmt"
	"log"
)

const (
	TYPE_NUMBER = "number"
	TYPE_BOOL = "bool"
)

type Variable struct {
	Name 	string			`data:"name"`
	vtype   string			`data:"value"`
	Value 	interface{}
	defined	bool // if the variable got a value or still not
}

func CreateNumber(name string) Variable {
	return Variable{
		Name: name,
		vtype: TYPE_NUMBER,
		defined: false,
	}
}

func CreateBool(name string) Variable {
	return Variable{
		Name: name,
		vtype: TYPE_BOOL,
		defined: false,
	}
}

func (v Variable) Clone() *Variable {
	return &Variable{
		Name: v.Name,
		Value: v.Value,
		vtype: v.vtype,
		defined: v.defined,
	}
}

func (v Variable) GetValue() (interface{}, bool) {
	if !v.defined {
		return nil, false
	}
	return v.Value, true
}

func (v *Variable) SetValue(vv interface{}) error {
	switch vv.(type) {
	case bool:
		if v.vtype == TYPE_BOOL {
			v.Value = vv
			v.defined = true
		} else {
			return fmt.Errorf("cannot set a bool to %v", v.Name)
		}
	case float64:
		if v.vtype == TYPE_NUMBER {
			v.Value = vv
			v.defined = true
		} else {
			return fmt.Errorf("cannot set a float64 to %v", v.Name)
		}
	default:
		return fmt.Errorf("unrecognized type of %v as value for %v", vv, v.Name)
	}
	return nil
}

func (v Variable) IsDefined() bool {
	return v.defined
}

func (v Variable) GetName() string {
	return v.Name
}

func (v Variable) Equals(vv Variable) bool {
	if v.Name != vv.Name || v.vtype != vv.vtype || v.IsDefined() != vv.IsDefined() {
		return false
	}
	if !v.IsDefined() {
		return true
	}
	v1, _ := v.GetValue()
	v2, _ := vv.GetValue()
	switch v.vtype {
	case TYPE_BOOL:
		if v1.(bool) != v2.(bool) {
			return false
		}
	case TYPE_NUMBER:
		if v1.(float64) != v2.(float64) {
			return false
		}
	default:
		log.Panicf("variable %v got an unrecognized type %v", v.Name, v.vtype)
	}
	return true
}