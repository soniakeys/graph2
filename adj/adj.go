// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

// Adj defines concrete types and methods for an adjacency graph.
// The types can be used by functions of graph/search for example.
package adj

import (
	"fmt"

	"github.com/soniakeys/graph"
)

// Node represents a node in an adjacency graph, either directed or undirected.
// It implements (for example) graph.AdjNode, graph.EstimateNode,
// graph.ArborNode, graph.SpannerNode and fmt.Stringer.
type Node struct {
	Data interface{}
	Nbs  []graph.Adj
}

// Visit visits adjacents of a Node.
func (n *Node) Visit(v graph.AdjVisitor) {
	for _, a := range n.Nbs {
		v(a)
	}
}

// Estimate obtains a heuristic distance estimate through the Estimator
// interface of the Data field of the receiver.  This panics if n.Data
// does not impliment Estimator.  It does not default to a null heuristic.
// Data is not required to implement Estimator, but if it does not implement
// Estimator, it is a programming error to use it in a context that requests
// estimates.
func (n *Node) Estimate(e graph.EstimateNode) float64 {
	return n.Data.(graph.Estimator).Estimate(e)
}

// String returns a string representation of n.Data.
func (n *Node) String() string { return fmt.Sprint(n.Data) }

// Weighted represents a weighted arc or edge.  It implements graph.Weighted.
type Weighted float64

// Weight returns arc or edge weight.
func (w Weighted) Weight() float64 { return float64(w) }

// Digraph defines a simple representation for a set of Nodes in a directed
// graph.
type Digraph map[interface{}]*Node

// Link sets one Node of a Digraph to be adjacent to another, adding either
// or both nodes to the graph as neccessary and adding an arc linking them.
//
// N1 and n2 are used as map keys and are also assigned to the Data fields
// when nodes are first added to the graph.  Because n1 and n2 are used as
// map keys, their concrete types must be Go-comparable.  N1 and n2 may be
// of any comparable type but if they will be used in a context that uses
// Nodes as graph.EstimateNodes for example, they must implement
// graph.Estimator.  Similarly, arc may be of any type but if Digraph g
// will be used in a context that uses arcs as weighted arcs, arc must
// implement graph.Weighted.
func (g Digraph) Link(n1, n2 interface{}, arc graph.Arc) {
	nd2, ok := g[n2]
	if !ok {
		nd2 = &Node{Data: n2}
		g[n2] = nd2
	}
	if nd1, ok := g[n1]; !ok {
		nd1 = &Node{n1, []graph.Adj{{arc, nd2}}}
		g[n1] = nd1
	} else {
		nd1.Nbs = append(nd1.Nbs, graph.Adj{arc, nd2})
	}
}

// LinkFrom lets Node satisfy graph.ArborNode, to enable creation of an
// arborescence on top of a graph.
func (n *Node) LinkFrom(prev graph.AdjNode, arc graph.Arc) graph.AdjNode {
	rn := &Node{Data: n} // create new node referring to receiver.
	if prev != nil {
		a := graph.Adj{Nd: rn}
		if wa, ok := arc.(graph.Weighted); ok {
			a.Ed = Weighted(wa.Weight()) // create arc if meaningful
		}
		pn := prev.(*Node)
		pn.Nbs = append(pn.Nbs, a)
	}
	return rn
}

// Span lets Node satisfy graph.SpannerNode, to enable creation of a
// spanning tree on top of a graph.
func (n *Node) Span(prev graph.AdjNode, ed graph.Edge) graph.AdjNode {
	rn := &Node{Data: n} // create new node referring to receiver.
	if prev != nil {
		a := graph.Adj{Nd: rn}
		if we, ok := ed.(graph.Weighted); ok {
			a.Ed = Weighted(we.Weight()) // create edge if meaningful
		}
		pn := prev.(*Node)
		pn.Nbs = append(pn.Nbs, a)
		// above code same as LinkFrom.  two lines below are new.
		a.Nd = prev
		rn.Nbs = []graph.Adj{a}
	}
	return rn
}

type Graph struct {
	Nodes map[interface{}]*Node
	Edges map[struct{ n1, n2 *Node }]graph.Edge
}

func NewGraph() Graph {
	return Graph{
		map[interface{}]*Node{},
		map[struct{ n1, n2 *Node }]graph.Edge{},
	}
}

func (g Graph) Link(n1, n2 interface{}, ed graph.Edge) {
	// get nodes first
	nd1, ok1 := g.Nodes[n1]
	if !ok1 {
		nd1 = &Node{Data: n1}
		g.Nodes[n1] = nd1
	}
	nd2, ok2 := g.Nodes[n2]
	if !ok2 {
		nd2 = &Node{Data: n2}
		g.Nodes[n2] = nd2
	}
	// if neither node was new, see if edge already exists
	if ok1 && ok2 {
		if _, ok := g.Edges[struct{ n1, n2 *Node }{nd1, nd2}]; ok {
			return // edge exists
		}
		if _, ok := g.Edges[struct{ n1, n2 *Node }{nd2, nd1}]; ok {
			return // edge exists
		}
	}
	// edge is new
	g.Edges[struct{ n1, n2 *Node }{nd1, nd2}] = ed
	nd1.Nbs = append(nd1.Nbs, graph.Adj{ed, nd2})
	nd2.Nbs = append(nd2.Nbs, graph.Adj{ed, nd1})
}
