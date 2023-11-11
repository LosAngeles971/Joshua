package graph

import "github.com/losangeles971/spiderweb/graph/core"

func (g Graph) GetConnectedToBy(node core.URI, predicate core.URI) ([]core.Relation, error) {
	return g.s.GetRelationsToNodeBy(node, predicate)
}