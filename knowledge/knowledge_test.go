package knowledge

import (
	"io/ioutil"
	"testing"
)

const (
	SUCCESS_EVENT = "They are all on the est bank of the river"
)
var	CAUSES = []string{
	"The farmer brings the cabbage to the est bank of the river",
	"The farmer brings the cabbage to the ovest bank of the river",
	"The farmer brings the goat to the est bank of the river",
	"The farmer brings the goat to the ovest bank of the river",
	"The farmer brings the wolf to the est bank of the river",
	"The farmer brings the wolf to the ovest bank of the river",
	"The farmer goes to the est bank of the river",
	"The farmer comes back to the ovest bank of the river",
}

func TestWhoCause(t *testing.T) {
	source, err := ioutil.ReadFile("../../../resources/the_farmer.joshua")
	if err != nil {
		t.Fatal(err)
	}
	k, err := Load(string(source))
	if err != nil {
		t.Fatal(err)
	}
	targetEvent, ok := k.GetEvent(SUCCESS_EVENT)
	if !ok {
		t.Fatalf("missing event %v", SUCCESS_EVENT)
	}
	causes := k.WhoCause(*targetEvent)
	if len(causes) != len(CAUSES) {
		t.Fatalf("expected %v causes not %v", len(CAUSES), len(causes))
	}
	for i := range CAUSES {
		ok := false
		for j := range causes {
			if causes[j].GetID() == CAUSES[i] {
				ok = true
			}
		}
		if !ok {
			t.Fatalf("missing cause %v", CAUSES[i])
		}
	}
}

func TestBranch(t *testing.T) {
	source, err := ioutil.ReadFile("../../../resources/testing_graph.joshua")
	if err != nil {
		t.Fatal(err)
	}
	k, err := Load(string(source))
	if err != nil {
		t.Fatal(err)
	}
	success, ok := k.GetEvent("Z")
	if !ok {
		t.Fatalf("missing event Z")
	}
	e_e, _ := k.GetEvent("E")
	e_f, _ := k.GetEvent("F")
	p1 := &Path{
		Path: []*Edge{
			&Edge{Cause: e_f, Effect: success,},
		},
	}
	edge_1 := Edge{
		Cause: e_e,
		Effect: e_f,
	} 
	p2 := getBranch(p1, &edge_1)
	if len(p1.Path) != 1 {
		t.Fatal("Corrupted p1")
	}
	if len(p2.Path) != 2 {
		t.Fatal("Corrupted p1")
	}
}

func TestAllPaths(t *testing.T) {
	source, err := ioutil.ReadFile("../../../resources/testing_graph.joshua")
	if err != nil {
		t.Fatal(err)
	}
	k, err := Load(string(source))
	if err != nil {
		t.Fatal(err)
	}
	e_z, ok := k.GetEvent("Z")
	if !ok {
		t.Fatal("missing event Z")
	}
	e_a, _ := k.GetEvent("A")
	e_b, _ := k.GetEvent("B")
	e_c, _ := k.GetEvent("C")
	e_d, _ := k.GetEvent("D")
	e_e, _ := k.GetEvent("E")
	e_f, _ := k.GetEvent("F")
	p1 := &Path{
		Path: []*Edge{
			&Edge{Cause: e_f, Effect: e_z,},
		},
	}
	p2 := &Path{
		Path: []*Edge{
			&Edge{Cause: e_c, Effect: e_z,},
			&Edge{Cause: e_b, Effect: e_c,},
			&Edge{Cause: e_a, Effect: e_b,},
			&Edge{Cause: e_e, Effect: e_a,},
		},
	}
	p3 := &Path{
		Path: []*Edge{
			&Edge{Cause: e_d, Effect: e_z,},
			&Edge{Cause: e_c, Effect: e_d,},
			&Edge{Cause: e_b, Effect: e_c,},
			&Edge{Cause: e_a, Effect: e_b,},
			&Edge{Cause: e_e, Effect: e_a,},
		},
	}
	allPaths := k.GetAllPathsToEvent(e_z)
	if len(allPaths) != 3 {
		t.Fatalf("expected %v paths not %v", 3, len(allPaths))
	}
	p1_c := 0
	p2_c := 0
	p3_c := 0
	for _, p := range allPaths {
		if p1.equals(p) {
			p1_c++
		}
		if p2.equals(p) {
			p2_c++
		}
		if p3.equals(p) {
			p3_c++
		}
	}
	if p1_c != 1 {
		t.Fatal("Path 1 not found F->Z")
	}
	if p2_c != 1 {
		t.Fatal("Path 2 not found C->Z")
	}
	if p3_c != 1 {
		t.Fatal("Path 3 not found C->D->Z")
	}
}