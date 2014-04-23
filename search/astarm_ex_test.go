// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package search_test

import (
	"fmt"

	"github.com/soniakeys/graph2"
	"github.com/soniakeys/graph2/search"
)

// AStarM requires a node type that implements graph2.EstimateNode and an
// edge type that implements graph2.Weighted.  Our two types:
type (
	monoNode struct {
		name string       // node name
		h    float64      // heuristic distance estimate to end node
		nbs  []graph2.Half // "neighbors," adjacent arcs and nodes
	}
	monoArc float64
)

// Two methods implement graph2.Estimator.
func (n *monoNode) VisitAdjHalfs(v graph2.AdjHalfVisitor) {
	for _, a := range n.nbs {
		v(a)
	}
}
func (n *monoNode) Estimate(graph2.EstimateNode) float64 { return n.h }

// One method implements graph2.Weighted.
func (a monoArc) Weight() float64 {
	return float64(a)
}

// Implement fmt.Stringer to make output easy.
func (n *monoNode) String() string { return n.name }

// One more method to make graph construction easy.
func (n *monoNode) link(n2 *monoNode, weight int) {
	n.nbs = append(n.nbs, graph2.Half{monoArc(weight), n2})
}

func ExampleAStarM() {
	a := &monoNode{name: "a", h: 19}
	b := &monoNode{name: "b", h: 20}
	c := &monoNode{name: "c", h: 10}
	d := &monoNode{name: "d", h: 6}
	e := &monoNode{name: "e", h: 0}
	f := &monoNode{name: "f", h: 9}
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

	p, l := search.AStarM(a, e)
	fmt.Println("Shortest path:", p)
	fmt.Println("Path length:", l)
	// Output:
	// Directed graph with 6 nodes, 9 edges
	// Shortest path: [{<nil> a} {9 c} {11 d} {6 e}]
	// Path length: 26
}
