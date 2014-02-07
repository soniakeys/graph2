// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package search_test

import (
	"fmt"

	"github.com/soniakeys/graph"
	"github.com/soniakeys/graph/search"
)

// AStarA requires a node type that implements graph.EstimateNode and an
// edge type that implements graph.Weighted.  Our two types:
type (
	estNode struct {
		name string       // node name
		h    float64      // heuristic distance estimate to end node
		nbs  []graph.Half // "neighbors," adjacent arcs and nodes
	}
	estArc float64
)

// Two methods implement graph.Estimator.
func (n *estNode) VisitAdj(v graph.HalfVisitor) {
	for _, a := range n.nbs {
		v(a)
	}
}
func (n *estNode) Estimate(graph.EstimateNode) float64 { return n.h }

// One method implements graph.Weighted.
func (a estArc) Weight() float64 {
	return float64(a)
}

// Implement fmt.Stringer to make output easy.
func (n *estNode) String() string { return n.name }

// One more method to make graph construction easy.
func (n *estNode) link(n2 *estNode, weight int) {
	n.nbs = append(n.nbs, graph.Half{estArc(weight), n2})
}

func ExampleAStarA() {
	a := &estNode{name: "a", h: 19}
	b := &estNode{name: "b", h: 20}
	c := &estNode{name: "c", h: 10}
	d := &estNode{name: "d", h: 6}
	e := &estNode{name: "e", h: 0}
	f := &estNode{name: "f", h: 9}
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

	p, l := search.AStarA(a, e)
	fmt.Println("Shortest path:", p)
	fmt.Println("Path length:", l)
	// Output:
	// Directed graph with 6 nodes, 9 edges
	// Shortest path: [{<nil> a} {9 c} {11 d} {6 e}]
	// Path length: 26
}
