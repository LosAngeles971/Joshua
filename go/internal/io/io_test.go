package io

import (
	"testing"
)

func TestIO(t *testing.T) {
	k, err := Load("../../../resources/k_contadino.yml")
	if len(k.Events) == 0 || err != nil {
		t.FailNow()
	}
}