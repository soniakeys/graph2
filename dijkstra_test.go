// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/soniakeys/graph"
)

// djNode represents a node in a directed graph.  It represents directed edges
// from the node with the handy DistanceNeighbor type from the graph package.
type djNode struct {
	nbs  []graph.DistanceNeighbor // directed edges as DistanceNeighbors
	name string                   // example application specific data
}

// djEdge is a simple number representing an edge length/distance/weight.
type djEdge float64

// djNode implements graph.DistanceNode, also fmt.Stringer
func (n *djNode) DistanceNeighbors([]graph.DistanceNeighbor) []graph.DistanceNeighbor {
	return n.nbs
}
func (n *djNode) String() string { return n.name }

// djEdge implements graph.DistanceEdge
func (e djEdge) Distance() float64 { return float64(e) }

var (
	djNodeData = []string{"a", "b", "c", "d", "e", "f"}
	djEdgeData = []struct {
		v1, v2 string
		l      float64
	}{
		{"a", "b", 7},
		{"a", "c", 9},
		{"a", "f", 14},
		{"b", "c", 10},
		{"b", "d", 15},
		{"c", "d", 11},
		{"c", "f", 2},
		{"d", "e", 6},
		{"e", "f", 9},
	}
)

// linkDjGraph constructs a linked representation of example data.
func linkDjGraph() (startNode, endNode *djNode) {
	all := map[string]*djNode{}
	for _, n := range djNodeData {
		all[n] = &djNode{name: n}
	}
	// link neighbors
	for _, ge := range djEdgeData {
		n1 := all[ge.v1]
		n1.nbs = append(n1.nbs,
			graph.DistanceNeighbor{djEdge(ge.l), all[ge.v2]})
	}
	return all["a"], all["e"]
}

func TestDijkstraDirected(t *testing.T) {
	// construct linked representation of example data
	startNode, endNode := linkDjGraph()
	// run Dijkstra's shortest path algorithm
	p, l := graph.DijkstraShortestPath(startNode, endNode)
	expected := "[{<nil> a} {9 c} {11 d} {6 e}]"
	got := fmt.Sprint(p)
	if got != expected {
		t.Fatal("wrong path.  expected", expected, "got", got)
	}
	if l != 26 {
		t.Fatal("wrong path length.  expected 26, got", l)
	}
	// reverse path should not exist with directed example data
	p, l = graph.DijkstraShortestPath(endNode, startNode)
	if p != nil {
		t.Fatal("wrong path.  expected nil, got", p)
	}
	if l != math.Inf(1) {
		t.Fatal("wrong path length.  expected +Inf, got", l)
	}
}
