package inmemory

import (
	"testing"

	"github.com/losangeles971/spiderweb/graph/core"
)

const (
	test_namespace = "test"
)

func TestLock(t *testing.T) {
	NewInMemoryStore()
	n1 := core.NewNode(test_namespace, "n1")
	n1.GetData().SetData("counter", int(0))
	for i :=1; i < 100; i++ {
	}
}