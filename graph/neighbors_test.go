package graph

import (
	"testing"

	"github.com/losangeles971/spiderweb/graph/core"
	"github.com/stretchr/testify/require"
)

func TestConnectedTo(t *testing.T) {
	g := NewInMemoryGraph()
	c1 := core.NewNode(test_namespace, "c1")
	n1 := core.NewNode(test_namespace, "n1")
	n2 := core.NewNode(test_namespace, "n2")
	n3 := core.NewNode(test_namespace, "n3")
	r1 := core.NewRelation(test_namespace, "to", n1.GetID(), c1.GetID())
	r2 := core.NewRelation(test_namespace, "to", n2.GetID(), c1.GetID())
	r3 := core.NewRelation(test_namespace, "to", c1.GetID(), n3.GetID())
	require.Nil(t, g.GetStore().StoreNode(c1))
	require.Nil(t, g.GetStore().StoreNode(n1))
	require.Nil(t, g.GetStore().StoreNode(n2))
	require.Nil(t, g.GetStore().StoreNode(n3))
	require.Nil(t, g.GetStore().StoreRelation(r1))
	require.Nil(t, g.GetStore().StoreRelation(r2))
	require.Nil(t, g.GetStore().StoreRelation(r3))
	rr, err := g.GetConnectedToBy(c1.GetID(), r1.GetPredicate())
	require.Nil(t, err)
	require.Equal(t, 2, len(rr))
}