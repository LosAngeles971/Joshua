package boltstore

/*++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

BoltStore is a persistent iplementation of a store interface, using BoltDB as back-end.

++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++*/

import (
	"fmt"
	"sync"

	"github.com/boltdb/bolt"
	"github.com/losangeles971/spiderweb/graph/core"
)

const (
	nodes_bucket     = "nodes"     // bucket dedicated to nodes
	relations_bucket = "relations" // bucket dedicated to relations
)

type BoltStore struct {
	path string
	db   *bolt.DB
	lock sync.Mutex
}

func NewBoltStore(path string) *BoltStore {
	return &BoltStore{
		path: path,
		db:   nil,
		lock: sync.Mutex{},
	}
}

// Open: it opens the database and creates the two buckets (nodes and relations) if they do not exist.
func (s *BoltStore) Open() error {
	if s.db == nil {
		var err error
		s.db, err = bolt.Open(s.path, 0600, nil)
		if err != nil {
			return err
		}
	}
	s.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(nodes_bucket))
		if err != nil {
			return fmt.Errorf("create bucket ( %s ) - %v", nodes_bucket, err)
		}
		return nil
	})
	s.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(relations_bucket))
		if err != nil {
			return fmt.Errorf("create bucket ( %s ) - %v", relations_bucket, err)
		}
		return nil
	})
	return nil
}

// Close: it closes the database.
func (s *BoltStore) Close() error {
	return s.db.Close()
}

// StoreNode: it stores or updates a node into the store.
func (s *BoltStore) StoreNode(n core.Node) error {
	dd, err := core.Encode(n)
	if err != nil {
		return err
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(nodes_bucket))
		err := b.Put([]byte(n.GetID().GetName()), dd)
		return err
	})
}

func (s *BoltStore) GetNodeByID(id core.URI) (core.Node, bool, error) {
	var node core.Node
	var ok bool
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(nodes_bucket))
		dd := b.Get([]byte(id.GetName()))
		if dd == nil {
			ok = false
			return nil
		}
		var err error
		node, err = core.DecodeNode(dd)
		if err != nil {
			return err
		}
		ok = true
		return nil
	})
	return node, ok, err
}

func (s *BoltStore) FindNodesByName(name core.URI) ([]core.Node, error) {
	result := []core.Node{}
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(nodes_bucket))
		c := b.Cursor()
		for k, dd := c.First(); k != nil; k, dd = c.Next() {
			node, err := core.DecodeNode(dd)
			if err != nil {
				return err
			}
			if node.Name.IsEqual(name) {
				result = append(result, node)
			}
		}
		return nil
	})
	return result, err
}

func (s *BoltStore) GetRelationByID(id core.URI) (core.Relation, bool, error) {
	var rel core.Relation
	var ok bool
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(relations_bucket))
		dd := b.Get([]byte(id.GetName()))
		if dd == nil {
			ok = false
			return nil
		}
		var err error
		rel, err = core.DecodeRelation(dd)
		if err != nil {
			return err
		}
		ok = true
		return nil
	})
	return rel, ok, err
}

func (s *BoltStore) StoreRelation(r core.Relation) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	dd, err := core.Encode(r)
	if err != nil {
		return err
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(relations_bucket))
		err := b.Put([]byte(r.GetID().GetName()), dd)
		return err
	})
}

func (s *BoltStore) findRelations(node core.URI, predicate core.URI, from bool, to bool, label bool) ([]core.Relation, error) {
	result := []core.Relation{}
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(relations_bucket))
		c := b.Cursor()
		for k, dd := c.First(); k != nil; k, dd = c.Next() {
			rel, err := core.DecodeRelation(dd)
			if err != nil {
				return err
			}
			if rel.GetSource().IsEqual(node) && from {
				if label {
					if rel.GetPredicate().IsEqual(predicate) {
						result = append(result, rel)
					}
				} else {
					result = append(result, rel)
				}
			}
			if rel.GetTarget().IsEqual(node) && to {
				result = append(result, rel)
			}
		}
		return nil
	})
	return result, err
}

func (s *BoltStore) GetRelationsFromNode(node core.URI) ([]core.Relation, error) {
	return s.findRelations(node, core.URI{}, true, false, false)
}

func (s *BoltStore) GetRelationsToNode(node core.URI) ([]core.Relation, error) {
	return s.findRelations(node, core.URI{}, false, true, false)
}

func (s *BoltStore) GetRelationsFromNodeBy(node core.URI, predicate core.URI) ([]core.Relation, error) {
	return s.findRelations(node, predicate, true, false, true)
}

func (s *BoltStore) GetRelationsToNodeBy(node core.URI, predicate core.URI) ([]core.Relation, error) {
	return s.findRelations(node, predicate, false, true, true)
}
