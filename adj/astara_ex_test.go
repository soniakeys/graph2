// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package adj_test

import (
	"fmt"

	"github.com/soniakeys/graph"
	"github.com/soniakeys/graph/adj"
	"github.com/soniakeys/graph/search"
)

// Adj types satisfy the interfaces required by AStarA.  The only additional
// requirement is a node data type that implements graph.Estimator.

// aaData represents node data.
type aaData struct {
	name string
	h    float64 // heuristic distance estimate to end node
}

// Implement graph.Estimator as required.
func (n *aaData) Estimate(graph.EstimateNode) float64 { return n.h }

// Implement fmt.String for nice output.
func (n *aaData) String() string { return n.name }

func ExampleDigraph_aStarA() {
	a := &aaData{"a", 19}
	b := &aaData{"b", 20}
	c := &aaData{"c", 10}
	d := &aaData{"d", 6}
	e := &aaData{"e", 0}
	f := &aaData{"f", 9}
	// construct graph using adj types
	g := adj.Digraph{}
	g.Link(a, b, adj.Weighted(7))
	g.Link(a, c, adj.Weighted(9))
	g.Link(a, f, adj.Weighted(14))
	g.Link(b, c, adj.Weighted(10))
	g.Link(b, d, adj.Weighted(15))
	g.Link(c, d, adj.Weighted(11))
	g.Link(c, f, adj.Weighted(2))
	g.Link(d, e, adj.Weighted(6))
	g.Link(e, f, adj.Weighted(9))
	// echo initial conditions
	fmt.Println("Directed graph with", len(g), "nodes")
	// run AStarA
	p, l := search.AStarA(g[a], g[e])
	fmt.Println("Shortest path:", p)
	fmt.Println("Path length:", l)
	// Output:
	// Directed graph with 6 nodes
	// Shortest path: [{<nil> a} {9 c} {11 d} {6 e}]
	// Path length: 26
}
