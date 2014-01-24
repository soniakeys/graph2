// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package search_test

import (
	"fmt"

	"github.com/soniakeys/graph"
	"github.com/soniakeys/graph/adj"
	"github.com/soniakeys/graph/search"
)

// amData represents data for a graph.EstimateNode.
type amData struct {
	name string  // example application specific data
	hEnd float64 // heuristic distance estimate to end node
}

// Implement graph.Estimator and fmt.String
func (n *amData) Estimate(graph.EstimateNode) float64 { return n.hEnd }
func (n *amData) String() string                      { return n.name }

func ExampleAStarM() {
	a := &amData{"a", 19}
	b := &amData{"b", 20}
	c := &amData{"c", 10}
	d := &amData{"d", 6}
	e := &amData{"e", 0}
	f := &amData{"f", 9}
	// use package graph/adj
	g := adj.Graph{}
	g.Link(a, b, adj.Edge(7))
	g.Link(a, c, adj.Edge(9))
	g.Link(a, f, adj.Edge(14))
	g.Link(b, c, adj.Edge(10))
	g.Link(b, d, adj.Edge(15))
	g.Link(c, d, adj.Edge(11))
	g.Link(c, f, adj.Edge(2))
	g.Link(d, e, adj.Edge(6))
	g.Link(e, f, adj.Edge(9))
	// echo initial conditions
	fmt.Println("Directed graph with", len(g), "nodes")
	// run AStarA
	p, l := search.AStarM(g[a], g[e])
	fmt.Println("Shortest path:", p)
	fmt.Println("Path length:", l)
	// Output:
	// Directed graph with 6 nodes
	// Shortest path: [{<nil> a} {9 c} {11 d} {6 e}]
	// Path length: 26
}
