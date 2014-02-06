// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package search_test

import (
	"fmt"

	"github.com/soniakeys/graph"
	"github.com/soniakeys/graph/search"
)

// DijkstraAllPaths requires a node type that implements graph.ArborNode
// and an edge type that implements graph.Weighted.  Our two types:

type (
	dapNode struct {
		name string       // node name
		nbs  []graph.Half // "neighbors," adjacent arcs and nodes
	}
	dapArc float64
)

// Two methods implement graph.ArborNode.
func (n *dapNode) Visit(v graph.HalfVisitor) {
	for _, a := range n.nbs {
		v(a)
	}
}
func (n *dapNode) LinkFrom(prev graph.AdjNode, arc graph.Arc) graph.AdjNode {
	rn := &arborNode{dap: n} // create new node referring to receiver.
	if prev != nil {
		a := graph.Half{Nd: rn}
		if wa, ok := arc.(graph.Weighted); ok {
			a.Ed = dapArc(wa.Weight()) // create arc if meaningful
		}
		pn := prev.(*arborNode)
		pn.nbs = append(pn.nbs, a)
	}
	return rn
}

// The node type for the arborescence can be merged with the node type of
// the underlying tree; see DijkstraAllPaths example in package graph/adj
// for example.  Here though we define a separate type to illustrate the
// separate roles.
type arborNode struct {
	dap *dapNode     // reference to the original graph
	nbs []graph.Half // branches from this node in the arboresence
}

// Satisfy graph.AdjNode.  (arborNode does not need to satisfy
// graph.ArborNode, the original graph nodes do.)
func (n *arborNode) Visit(v graph.HalfVisitor) {
	for _, a := range n.nbs {
		v(a)
	}
}

// Implement fmt.Stringer on arborNode.  These are the nodes we will output.
func (n *arborNode) String() string {
	return n.dap.name // pull node name from original graph node
}

// One method implements graph.Weighted.
func (a dapArc) Weight() float64 {
	return float64(a)
}

// One more method on dapNode to make graph construction easy.
func (n *dapNode) link(n2 *dapNode, weight int) {
	n.nbs = append(n.nbs, graph.Half{dapArc(weight), n2})
}

func ExampleDijkstraAllPaths() {
	a := &dapNode{name: "a"}
	b := &dapNode{name: "b"}
	c := &dapNode{name: "c"}
	d := &dapNode{name: "d"}
	e := &dapNode{name: "e"}
	f := &dapNode{name: "f"}
	a.link(b, 7)
	a.link(c, 9)
	a.link(f, 14)
	b.link(c, 10)
	b.link(d, 15)
	c.link(d, 11)
	c.link(f, 2)
	d.link(e, 6)
	e.link(f, 9)
	// a recursive function to print paths
	var pp func(string, graph.AdjNode)
	pp = func(s string, n graph.AdjNode) {
		s += fmt.Sprint(n)
		fmt.Println(s)
		s += " "
		n.Visit(func(nb graph.Half) {
			pp(s, nb.Nd)
		})
	}
	// run Dijkstra's algorithm to find all shortest paths
	pp("", search.DijkstraAllPaths(a))
	// Output:
	// a
	// a b
	// a c
	// a c f
	// a c d
	// a c d e
}
