// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package search_test

import (
	"fmt"

	"github.com/soniakeys/graph"
	"github.com/soniakeys/graph/adj"
	"github.com/soniakeys/graph/search"
)

// aaData represents data for a graph.EstimateNode.
type aaData struct {
	name string  // example application specific data
	hEnd float64 // heuristic distance estimate to end node
}

// Implement graph.Estimator and fmt.String
func (n *aaData) Estimate(graph.EstimateNode) float64 { return n.hEnd }
func (n *aaData) String() string                      { return n.name }

func ExampleAStarA() {
	a := &aaData{"a", 19}
	b := &aaData{"b", 20}
	c := &aaData{"c", 10}
	d := &aaData{"d", 6}
	e := &aaData{"e", 0}
	f := &aaData{"f", 9}
	// use package graph/adj
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
