package problems

import (
	"testing"
	"fmt"
)

func TestLoadProblem(t *testing.T) {
	init, success_name, err := Load("../../../resources/p_contadino.yml")
	if err != nil {
		fmt.Println("Error loading the problem: ", err)
		t.FailNow()
	}
	if init.Size() != 8 {
		fmt.Println("Wrong number of variables: ", init.Size())
		t.FailNow()
	}
	if success_name != "The farmer, the wolf, the goat and the cabbage are on the bank B of the river" {
		fmt.Println("Success not parsed")
		t.FailNow()
	}
}

func TestLoadProblem2(t *testing.T) {
	init, success_name, err := Load("../../../resources/p_aereo.yml")
	if err != nil {
		fmt.Println("Error loading the problem: ", err)
		t.FailNow()
	}
	if init.Size() != 14 {
		fmt.Println("Wrong number of variables: ", init.Size())
		t.FailNow()
	}
	if success_name != "I know their jobs" {
		fmt.Println("Success not parsed")
		t.FailNow()
	}
}