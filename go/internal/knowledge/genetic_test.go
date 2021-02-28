package knowledge

import (
	"fmt"
	ctx "it/losangeles971/joshua/internal/context"
	"testing"
)

func TestGenetic(t *testing.T) {
	init := ctx.CreateEmptyState()
	ra, err := ctx.ParseRange("{1,2}")
	if err != nil {
		t.FailNow()
	}
	rb, err := ctx.ParseRange("{3,4,5}")
	if err != nil {
		t.FailNow()
	}
	rc, err := ctx.ParseRange("{1,2,3,4,5}")
	if err != nil {
		t.FailNow()
	}
	init.Add(&ctx.Variable{Name: "A", Defined: false, Range: ra})
	init.Add(&ctx.Variable{Name: "B", Defined: false, Range: rb})
	init.Add(&ctx.Variable{Name: "C", Defined: false, Range: rc})
	init.Add(&ctx.Variable{Name: "D", Defined: true, Value: 0})
	e1 := Edge{
		Cause: &Event{
			ID: "1",
			Conditions: []string{
				"A == 2",
				"B == 5",
			},
		},
		Effect: &Event{
			ID: "2",
			Conditions: []string{
				"C > 3",
				"D == 0",
			},
		},
	}
	p1 := Path{Path: []*Edge{&e1}}
	p := MakePopulation(&p1, init)
	err = CycleGenerations(&p, 10)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	if len(p.Population) != 10 {
		t.FailNow()
	}
	solutions := 0
	for i := range p.Population {
		a, ok := p.Population[i].DNA["A"]
		if !ok {
			t.FailNow()
		}
		b, ok := p.Population[i].DNA["B"]
		if !ok {
			t.FailNow()
		}
		c, ok := p.Population[i].DNA["C"]
		if !ok {
			t.FailNow()
		}
		d, ok := p.Population[i].DNA["D"]
		if !ok {
			t.FailNow()
		}
		if a == 2 && b == 5 && c > 3 {
			solutions++
		}
		fmt.Println("Person ", a, b, c, d, p.Population[i].ranking)
	}
	fmt.Println("Solutions found ", solutions)
	if solutions < 1 {
		fmt.Println("No solutions")
		t.FailNow()
	}
	ok, err := p.GetOneSolution(&init)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	if !ok {
		fmt.Println("No solutions")
		t.FailNow()
	}
	p1.Reset()
	err = p1.Run(init, 0)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	if p1.Outcome != CE_OUTCOME_TRUE {
		fmt.Println("Corrupted solution")
		t.FailNow()
	}
}