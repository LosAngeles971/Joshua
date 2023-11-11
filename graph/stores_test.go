package graph

import (
	"os"
	"testing"

	"github.com/losangeles971/spiderweb/graph/boltstore"
	"github.com/losangeles971/spiderweb/graph/core"
	"github.com/losangeles971/spiderweb/graph/inmemory"
	"github.com/stretchr/testify/require"
)

const (
	test_namespace = "test"
	test_predicate = "to"
)

var test_data = map[string]interface{}{
	"i":   int(15),
	"i64": int64(1500),
	"f":   float64(6.3),
	"s":   "test",
	"b":   false,
}

func testBasicStore(t *testing.T, s Store) {
	n1 := core.NewNode(test_namespace, "n1")
	require.NotNil(t, n1.GetData())
	err := n1.GetData().SetDatas(test_data)
	require.Nil(t, err)
	err = s.StoreNode(n1)
	require.Nil(t, err)
	n2, ok, err := s.GetNodeByID(n1.GetID())
	require.Nil(t, err)
	require.True(t, ok)
	require.True(t, n1.IsEqual(n2))
	r1 := core.NewRelation(test_namespace, test_predicate, n1.ID, n1.ID)
	require.Nil(t, err)
	err = s.StoreRelation(r1)
	require.Nil(t, err)
	r2, ok, err := s.GetRelationByID(r1.GetID())
	require.Nil(t, err)
	require.True(t, ok)
	require.True(t, r1.IsEqual(r2))
}

func TestBasicStoreInMemory(t *testing.T) {
	s := inmemory.NewInMemoryStore()
	testBasicStore(t, s)
}

func TestBasicStoreBolt(t *testing.T) {
	tmpdir, err := os.MkdirTemp(os.TempDir(), "spiderweb")
	require.Nil(t, err)
	s := boltstore.NewBoltStore(tmpdir + "/test")
	require.Nil(t, s.Open())
	testBasicStore(t, s)
	require.Nil(t, s.Close())	
}

func testGetRelationsToNode(t *testing.T, s Store) {
	c1 := core.NewNode(test_namespace, "c1")
	n1 := core.NewNode(test_namespace, "n1")
	n2 := core.NewNode(test_namespace, "n2")
	n3 := core.NewNode(test_namespace, "n3")
	r1 := core.NewRelation(test_namespace, "to", n1.GetID(), c1.GetID())
	r2 := core.NewRelation(test_namespace, "to", n2.GetID(), c1.GetID())
	r3 := core.NewRelation(test_namespace, "to", c1.GetID(), n3.GetID())
	require.Nil(t, s.StoreNode(c1))
	require.Nil(t, s.StoreNode(n1))
	require.Nil(t, s.StoreNode(n2))
	require.Nil(t, s.StoreNode(n3))
	require.Nil(t, s.StoreRelation(r1))
	require.Nil(t, s.StoreRelation(r2))
	require.Nil(t, s.StoreRelation(r3))
	to, err := s.GetRelationsToNodeBy(c1.GetID(), r1.GetPredicate())
	require.Nil(t, err)
	require.Equal(t, 2, len(to))
	from, err := s.GetRelationsFromNodeBy(c1.GetID(), r1.GetPredicate())
	require.Nil(t, err)
	require.Equal(t, 1, len(from))
}

func TestGetRelationsToNodeInMemory(t *testing.T) {
	s := inmemory.NewInMemoryStore()
	testGetRelationsToNode(t, s)
}

func TestGetRelationsToNodeBolt(t *testing.T) {
	tmpdir, err := os.MkdirTemp(os.TempDir(), "spiderweb")
	require.Nil(t, err)
	s := boltstore.NewBoltStore(tmpdir + "/test")
	require.Nil(t, s.Open())
	testGetRelationsToNode(t, s)
	require.Nil(t, s.Close())	
}