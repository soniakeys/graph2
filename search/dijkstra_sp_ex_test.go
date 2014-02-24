// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package search_test

import (
	"fmt"

	"github.com/soniakeys/graph"
	"github.com/soniakeys/graph/search"
)

// DijkstraShortestPath requires a node type that implements graph.AdjNode
// and an edge type that implements graph.Weighted.  Our two types:

type (
	dspNode struct {
		name string       // node name
		nbs  []graph.Half // "neighbors," adjacent arcs and nodes
	}
	dspArc float64
)

// One method implements graph.AdjNode.
func (n *dspNode) VisitAdjHalfs(v graph.AdjHalfVisitor) {
	for _, a := range n.nbs {
		v(a)
	}
}

// One method implements graph.Weighted.
func (a dspArc) Weight() float64 {
	return float64(a)
}

// Implement fmt.Stringer to make output easy.
func (n *dspNode) String() string { return n.name }

// One more method to make graph construction easy.
func (n *dspNode) link(n2 *dspNode, weight int) {
	n.nbs = append(n.nbs, graph.Half{dspArc(weight), n2})
}

func ExampleDijkstraShortestPath_directed() {
	a := &dspNode{name: "a"}
	b := &dspNode{name: "b"}
	c := &dspNode{name: "c"}
	d := &dspNode{name: "d"}
	e := &dspNode{name: "e"}
	f := &dspNode{name: "f"}
	a.link(b, 7)
	a.link(c, 9)
	a.link(f, 14)
	b.link(c, 10)
	b.link(d, 15)
	c.link(d, 11)
	c.link(f, 2)
	d.link(e, 6)
	e.link(f, 9)
	fmt.Println("Directed graph with 6 nodes, 9 edges")

	path, l := search.DijkstraShortestPath(a, e)
	fmt.Println(`Shortest path from node "a" to node "e":`, path)
	fmt.Println("Path length:", l)
	// Output:
	// Directed graph with 6 nodes, 9 edges
	// Shortest path from node "a" to node "e": [{<nil> a} {9 c} {11 d} {6 e}]
	// Path length: 26
}
