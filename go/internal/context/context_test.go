package context

import (
	"testing"
	"fmt"
)

func TestPartOf(t *testing.T) {
	c1 := Create()
	c2 := Create()
	c1.Add(&Variable{Name: "a", Value: 1.0, Defined: true, })
	c1.Add(&Variable{Name: "b", Value: -1.0, Defined: true, })
	c2.Add(&Variable{Name: "a", Value: 1.0, Defined: true, })
	c2.Add(&Variable{Name: "b", Value: -1.0, Defined: true, })
	c2.Add(&Variable{Name: "c", Value: 0.0, Defined: true, })
	if ok := c1.PartOf(c2); !ok {
		fmt.Println("c1 must be part of c2!")
		t.FailNow()
	}
	if ok := c2.PartOf(c1); ok {
		t.FailNow()
	}
	c3 := c2.Clone()
	c3.Add(&Variable{Name: "c", Value: 0.4, Defined: true, })
	if ok := c2.PartOf(c3); ok {
		t.FailNow()
	}
}

func TestRange(t *testing.T) {
	c1 := Create()
	r, err := ParseRange("{1, 2, 3}")
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	c1.Add(&Variable{Name: "a", Range: r, })
	if c1.data["a"].Range.Size() != 3 {
		t.FailNow()
	}
}