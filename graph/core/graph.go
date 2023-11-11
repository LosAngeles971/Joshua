package core

import (
	"bytes"
	"encoding/gob"
)

type Node struct {
	ID   URI
	Name URI
	Data *Data
}

type Relation struct {
	ID        URI
	Predicate URI
	Data      *Data
	Source    URI
	Target    URI
}

func NewNode(namespace string, name string) Node {
	return Node{
		ID:   NewID(namespace),
		Name: NewURI(namespace, name),
		Data: newData(),
	}
}

func NewRelation(namespace string, predicate string, source URI, target URI) Relation {
	return Relation{
		ID:        NewID(namespace),
		Predicate: NewURI(namespace, predicate),
		Source:    source,
		Target:    target,
		Data:      newData(),
	}
}

func (n Node) GetID() URI {
	return n.ID
}

func (n Node) GetName() URI {
	return n.Name
}

func (n Node) GetData() *Data {
	return n.Data
}

func (n Node) Clone() Node {
	return Node{
		ID:   n.ID.Clone(),
		Name: n.Name.Clone(),
		Data: n.Data.Clone(),
	}
}

func (n1 Node) IsEqual(n2 Node) bool {
	if !n1.ID.IsEqual(n2.ID) {
		return false
	}
	if !n1.Name.IsEqual(n2.Name) {
		return false
	}
	if !n1.Data.IsEqual(n2.Data) {
		return false
	}
	return true
}

func Encode(obj interface{}) ([]byte, error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(obj)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func DecodeNode(data []byte) (Node, error) {
	node := Node{}
	buf := bytes.Buffer{}
	_, err := buf.Write(data)
	if err != nil {
		return node, err
	}
    dec := gob.NewDecoder(&buf)
    err = dec.Decode(&node)
    return node, err
}

func DecodeRelation(data []byte) (Relation, error) {
	rel := Relation{}
	buf := bytes.Buffer{}
	_, err := buf.Write(data)
	if err != nil {
		return rel, err
	}
    dec := gob.NewDecoder(&buf)
    err = dec.Decode(&rel)
    return rel, err
}

func (r Relation) GetID() URI {
	return r.ID
}

func (r Relation) GetPredicate() URI {
	return r.Predicate
}

func (r Relation) GetSource() URI {
	return r.Source
}

func (r Relation) GetTarget() URI {
	return r.Target
}

func (r Relation) GetData() *Data {
	return r.Data
}

func (r Relation) Clone() Relation {
	return Relation{
		ID:        r.ID.Clone(),
		Predicate: r.Predicate.Clone(),
		Source:    r.Source.Clone(),
		Target:    r.Target.Clone(),
		Data:      r.Data.Clone(),
	}
}

func (r1 Relation) IsEqual(r2 Relation) bool {
	if !r1.ID.IsEqual(r2.ID) {
		return false
	}
	if !r1.Predicate.IsEqual(r2.Predicate) {
		return false
	}
	if !r1.Data.IsEqual(r2.Data) {
		return false
	}
	if !r1.Source.IsEqual(r2.Source) {
		return false
	}
	if !r1.Target.IsEqual(r2.Target) {
		return false
	}
	return true
}