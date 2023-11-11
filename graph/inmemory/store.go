package inmemory

/*++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

InMemoryStore is a minimal implementation of a store interface, using memory as back-end.
InMemoryStore is useful for supporting the development and test of program which
use SpiderWeb library.
InMemoryStore is not designed for performance and production-like environment.

++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++*/

import (
	"fmt"
	"sync"

	"github.com/losangeles971/spiderweb/graph/core"
)

type InMemoryStore struct {
	nodesSize     int // maximum amount of nodes to keep in-memory
	relationsSize int // maximum amount of relations to keep in-memory
	// partioning between nodes and relations just to enhanche performance a little bit
	nodes     map[string]core.Node
	relations map[string]core.Relation
	lock      sync.Mutex
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		nodesSize:     1000,
		relationsSize: 1000,
		nodes:         map[string]core.Node{},
		relations:     map[string]core.Relation{},
		lock:          sync.Mutex{},
	}
}

func (s *InMemoryStore) Open() error {
	return nil
}

func (s *InMemoryStore) Close() error {
	return nil
}

func (s *InMemoryStore) GetNodeByID(id core.URI) (core.Node, bool, error) {
	if n, ok := s.nodes[id.GetName()]; ok {
		c := n.Clone()
		return c, true, nil
	}
	return core.Node{}, false, nil
}

func (s *InMemoryStore) FindNodesByName(name core.URI) ([]core.Node, error) {
	result := []core.Node{}
	for _, node := range s.nodes {
		if node.GetName().IsEqual(name) {
			result = append(result, node)
		}
	}
	return result, nil
}

func (s *InMemoryStore) StoreNode(n core.Node) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.nodes[n.GetID().GetName()]; ok {
		s.nodes[n.GetID().GetName()] = n.Clone()
	} else {
		if len(s.nodes) >= s.nodesSize {
			return fmt.Errorf("cannot store new node ( %s ) reached maximum amount of nodes ( %d )", n.GetID().GetName(), s.nodesSize)
		}
		s.nodes[n.GetID().GetName()] = n.Clone()
	}
	return nil
}

func (s *InMemoryStore) GetRelationByID(id core.URI) (core.Relation, bool, error) {
	if r, ok := s.relations[id.GetName()]; ok {
		c := r.Clone()
		return c, true, nil
	}
	return core.Relation{}, false, nil
}

func (s *InMemoryStore) StoreRelation(r core.Relation) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.relations[r.GetID().GetName()]; ok {
		s.relations[r.GetID().GetName()] = r.Clone()
	} else {
		if len(s.relations) >= s.relationsSize {
			return fmt.Errorf("cannot store new relation ( %s ) reached maximum amount of relations ( %d )", r.GetID().GetName(), s.relationsSize)
		}
		s.relations[r.GetID().GetName()] = r.Clone()
	}
	return nil
}

func (s *InMemoryStore) findRelations(node core.URI, predicate core.URI, from bool, to bool, label bool) ([]core.Relation, error) {
	result := []core.Relation{}
	for _, r := range s.relations {
		if r.GetSource().IsEqual(node) && from {
			if label {
				if r.GetPredicate().IsEqual(predicate) {
					result = append(result, r)
				}
			} else {
				result = append(result, r)
			}
		}
		if r.GetTarget().IsEqual(node) && to {
			result = append(result, r)
		}
	}
	return result, nil
}

func (s *InMemoryStore) GetRelationsFromNode(node core.URI) ([]core.Relation, error) {
	return s.findRelations(node, core.URI{}, true, false, false)
}

func (s *InMemoryStore) GetRelationsToNode(node core.URI) ([]core.Relation, error) {
	return s.findRelations(node, core.URI{}, false, true, false)
}

func (s *InMemoryStore) GetRelationsFromNodeBy(node core.URI, predicate core.URI) ([]core.Relation, error) {
	return s.findRelations(node, predicate, true, false, true)
}

func (s *InMemoryStore) GetRelationsToNodeBy(node core.URI, predicate core.URI) ([]core.Relation, error) {
	return s.findRelations(node, predicate, false, true, true)
}