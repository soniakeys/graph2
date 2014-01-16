// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

import (
	"math"
	"testing"

	"github.com/soniakeys/graph"
)

// minimal node type with no neighbors.
type djNode0 struct{}

func (n *djNode0) DistanceNeighbors([]graph.DistanceNeighbor) []graph.DistanceNeighbor {
	return nil
}

func TestDijkstraDirected(t *testing.T) {
	// search from node with no neighbors to something other than itself
	// should return no path.
	p, l := graph.DijkstraShortestPath(&djNode0{}, nil)
	if p != nil {
		t.Fatal("wrong path.  expected nil, got", p)
	}
	if l != math.Inf(1) {
		t.Fatal("wrong path length.  expected +Inf, got", l)
	}
}
