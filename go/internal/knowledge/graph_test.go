package knowledge

import (
	"fmt"
	"testing"
)

func TestStack(t *testing.T) {
	k := Knowledge{}
	err := k.Load("../../../resources/k_graph.yml")
	fmt.Println("Events: ", len(k.Events))
	if err != nil {
		fmt.Println("Error: ", err)
		t.FailNow()
	}
	success, ok := k.GetEvent("Z")
	if !ok {
		fmt.Println("Missing success")
		t.FailNow()
	}
	e_a, _ := k.GetEvent("A")
	e_b, _ := k.GetEvent("B")
	e_c, _ := k.GetEvent("C")
	e_d, _ := k.GetEvent("D")
	e_e, _ := k.GetEvent("E")
	e_f, _ := k.GetEvent("F")
	p1 := &Path{
		Path: []Edge{
			Edge{Cause: e_f, Effect: success,},
		},
	}
	p2 := &Path{
		Path: []Edge{
			Edge{Cause: e_d, Effect: e_e,},
		},
	}
	p3 := &Path{
		Path: []Edge{
			Edge{Cause: e_c, Effect: e_d,},
		},
	}
	p4 := &Path{
		Path: []Edge{
			Edge{Cause: e_b, Effect: e_c,},
			Edge{Cause: e_a, Effect: e_b,},
		},
	}
	ss := []*Path{p1, p2, p3, p4}
	s := Stack{}
	s.Push(p1)
	s.Push(p2)
	s.Push(p3)
	s.Push(p4)
	if len(ss) != 4 || s.Size() != 4 {
		t.FailNow()
	}
	s.Pop()
	s.Pop()
	s.Pop()
	s.Pop()
	if len(ss) != 4 || s.Size() != 0 {
		t.FailNow()
	}
}

func TestBranch(t *testing.T) {
	k := Knowledge{}
	err := k.Load("../../../resources/k_graph.yml")
	fmt.Println("Events: ", len(k.Events))
	if err != nil {
		fmt.Println("Error: ", err)
		t.FailNow()
	}
	success, ok := k.GetEvent("Z")
	if !ok {
		fmt.Println("Missing success")
		t.FailNow()
	}
	e_e, _ := k.GetEvent("E")
	e_f, _ := k.GetEvent("F")
	p1 := &Path{
		Path: []Edge{
			Edge{Cause: e_f, Effect: success,},
		},
	}
	edge_1 := Edge{
		Cause: e_e,
		Effect: e_f,
	} 
	p2 := getBranch(p1, edge_1)
	if len(p1.Path) != 1 {
		fmt.Println("Corrupted p1")
		t.FailNow()
	}
	if len(p2.Path) != 2 {
		fmt.Println("Corrupted p1")
		t.FailNow()
	}
}

func TestAllPaths(t *testing.T) {
	k := Knowledge{}
	err := k.Load("../../../resources/k_graph.yml")
	fmt.Println("Events: ", len(k.Events))
	if err != nil {
		fmt.Println("Error: ", err)
		t.FailNow()
	}
	e_z, ok := k.GetEvent("Z")
	if !ok {
		fmt.Println("Missing success")
		t.FailNow()
	}
	e_a, _ := k.GetEvent("A")
	e_b, _ := k.GetEvent("B")
	e_c, _ := k.GetEvent("C")
	e_d, _ := k.GetEvent("D")
	e_e, _ := k.GetEvent("E")
	e_f, _ := k.GetEvent("F")
	p1 := &Path{
		Path: []Edge{
			Edge{Cause: e_f, Effect: e_z,},
		},
	}
	p2 := &Path{
		Path: []Edge{
			Edge{Cause: e_c, Effect: e_z,},
			Edge{Cause: e_b, Effect: e_c,},
			Edge{Cause: e_a, Effect: e_b,},
			Edge{Cause: e_e, Effect: e_a,},
		},
	}
	p3 := &Path{
		Path: []Edge{
			Edge{Cause: e_d, Effect: e_z,},
			Edge{Cause: e_c, Effect: e_d,},
			Edge{Cause: e_b, Effect: e_c,},
			Edge{Cause: e_a, Effect: e_b,},
			Edge{Cause: e_e, Effect: e_a,},
		},
	}
	allPaths := GetAllPaths(k, e_z)
	if len(allPaths) != 3 {
		fmt.Println("No correct number of paths, expected 3 but: ", len(allPaths))
		t.FailNow()
	}
	p1_c := 0
	p2_c := 0
	p3_c := 0
	for _, p := range allPaths {
		if p1.Equals(p) {
			p1_c++
		}
		if p2.Equals(p) {
			p2_c++
		}
		if p3.Equals(p) {
			p3_c++
		}
	}
	if p1_c != 1 {
		fmt.Println("Path 1 not found F->Z")
		t.FailNow()		
	}
	if p2_c != 1 {
		fmt.Println("Path 2 not found C->Z")
		t.FailNow()		
	}
	if p3_c != 1 {
		fmt.Println("Path 3 not found C->D->Z")
		t.FailNow()		
	}
}