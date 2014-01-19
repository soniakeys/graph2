// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package search_test

import (
	"math"
	"testing"

	"github.com/soniakeys/graph"
	"github.com/soniakeys/graph/search"
)

// minimal node type with no neighbors.
type djNode0 struct{}

func (n *djNode0) Visit(graph.NeighborVisitor) {
}

func TestDijkstraDirected(t *testing.T) {
	// search from node with no neighbors to something other than itself
	// should return no path.
	np, ep, l := search.DijkstraShortestPath(&djNode0{}, nil)
	if np != nil {
		t.Fatal("wrong node path.  expected nil, got", np)
	}
	if ep != nil {
		t.Fatal("wrong edge path.  expected nil, got", ep)
	}
	if l != math.Inf(1) {
		t.Fatal("wrong path length.  expected +Inf, got", l)
	}
}
