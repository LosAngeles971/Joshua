package graph

/*++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

This package implements the factory design pattern.
The Store interface defines the mandatory functions of a generic implementation of the Store.
The package provides one method for each supported implementation of the Store interface.

++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++*/

import (
	"github.com/losangeles971/spiderweb/graph/boltstore"
	"github.com/losangeles971/spiderweb/graph/core"
	"github.com/losangeles971/spiderweb/graph/inmemory"
)

type Store interface {
	Open() error
	Close() error
	GetNodeByID(id core.URI) (core.Node, bool, error)
	FindNodesByName(name core.URI) ([]core.Node, error)
	StoreNode(n core.Node) error
	GetRelationByID(id core.URI) (core.Relation, bool, error)
	StoreRelation(n core.Relation) error
	GetRelationsFromNode(id core.URI) ([]core.Relation, error)
	GetRelationsToNode(id core.URI) ([]core.Relation, error)
	GetRelationsFromNodeBy(id core.URI, predicate core.URI) ([]core.Relation, error)
	GetRelationsToNodeBy(id core.URI, predicate core.URI) ([]core.Relation, error)
}

type Graph struct {
	s Store
}

func NewInMemoryGraph() Graph {
	return Graph{
		s: inmemory.NewInMemoryStore(),
	}
}

func NewBoltGraph(pathfile string) Graph {
	return Graph{
		s: boltstore.NewBoltStore(pathfile),
	}
}

func (g Graph) GetStore() Store {
	return g.s
}