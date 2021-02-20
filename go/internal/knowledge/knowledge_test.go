package knowledge

import (
	"fmt"
	"testing"
)

func TestKnowledge(t *testing.T) {
	k := Knowledge{}
	err := k.Load("../../../resources/k_contadino.yml")
	fmt.Println("Events: ", len(k.Events))
	if len(k.Events) != 9 || err != nil {
		fmt.Println("Error: ", err)
		t.FailNow()
	}
	success, ok := k.GetEvent("The farmer, the wolf, the goat and the cabbage are on the bank B of the river")
	if !ok {
		fmt.Println("Missing success")
		t.FailNow()
	}
	ee := k.IsEffectOf(success)
	if len(ee) != 8 {
		fmt.Println("Corrupted knowledge")
		t.FailNow()
	}
}
