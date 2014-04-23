// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package search_test

import (
	"fmt"

	"github.com/soniakeys/graph2"
	"github.com/soniakeys/graph2/search"
)

// DijkstraShortestPath requires a node type that implements graph2.AdjNode
// and an edge type that implements graph2.Weighted.  Our two types look
// the same as for the directed example but we will implement edges with
// reciprocal half arcs referencing a common edge object.

type (
	uNode struct {
		name string       // node name
		nbs  []graph2.Half // "neighbors," adjacent arcs and nodes
	}
	uEdge float64
)

// One method implements graph2.AdjNode.
func (n *uNode) VisitAdjHalfs(v graph2.AdjHalfVisitor) {
	for _, a := range n.nbs {
		v(a)
	}
}

// One method implements graph2.Weighted.
func (e uEdge) Weight() float64 {
	return float64(e)
}

// Implement fmt.Stringer on both node and edge types.
func (n *uNode) String() string { return n.name }
func (e uEdge) String() string  { return fmt.Sprint(float64(e)) }

// Method to make graph construction easy.
func (n1 *uNode) link(n2 *uNode, weight int) {
	e := uEdge(weight)
	n1.nbs = append(n1.nbs, graph2.Half{&e, n2})
	n2.nbs = append(n2.nbs, graph2.Half{&e, n1})
}

func ExampleDijkstraShortestPath_undirected() {
	a := &uNode{name: "a"}
	b := &uNode{name: "b"}
	c := &uNode{name: "c"}
	d := &uNode{name: "d"}
	e := &uNode{name: "e"}
	f := &uNode{name: "f"}
	a.link(b, 7)
	a.link(c, 9)
	a.link(f, 14)
	b.link(c, 10)
	b.link(d, 15)
	c.link(d, 11)
	c.link(f, 2)
	d.link(e, 6)
	e.link(f, 9)
	fmt.Println("Undirected graph with 6 nodes, 9 edges")

	path, l := search.DijkstraShortestPath(a, e)
	fmt.Println("Shortest path:", path)
	fmt.Println("Path length:", l)
	// Output:
	// Undirected graph with 6 nodes, 9 edges
	// Shortest path: [{<nil> a} {9 c} {2 f} {9 e}]
	// Path length: 20
}
