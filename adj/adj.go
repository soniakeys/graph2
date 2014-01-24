// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

// Adj defines concrete types and methods for an adjacency graph.
// The types can be used by functions of graph/search for example.
package adj

import (
	"fmt"

	"github.com/soniakeys/graph"
)

// Node represents a node in an adjacency graph.  It implements
// graph.EstimateNode and fmt.Stringer.
type Node struct {
	Data interface{}      // application specific
	Nbs  []graph.Neighbor // directed edges
}

// Visit visits neighbors of a node.
func (n *Node) Visit(v graph.NeighborVisitor) {
	for _, nb := range n.Nbs {
		v(nb)
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

// Edge implements graph.DistanceEdge.
type Edge float64

// Distance returns edge distance.
func (e Edge) Distance() float64 { return float64(e) }

// String returns a string representation of the edge distance.
func (e Edge) String() string { return fmt.Sprint(float64(e)) }

// Graph defines a simple representation for a set of Nodes.
type Graph map[interface{}]*Node

// Link sets one Node of a Graph to be a neighbor of another, adding either
// or both nodes to the graph as neccessary.
//
// N1 and n2 are used as node keys and are also assigned to the Data fields
// when nodes are first added to the graph.  Because n1 and n2 are used as
// map keys, their concrete types must be Go-comparable.  N1 and n2 may be
// of any comparable type but if they will be used in a context that uses
// Nodes as graph.EstimateNodes for example, they must implement
// graph.Estimator.  Similarly, ed may be of any type, but it will be used
// as a graph.Edge.  If the Graph g will be used in a context that uses edges
// as graph.DistanceEdges, ed must implement graph.DistanceEdge.
func (g Graph) Link(n1, n2, ed interface{}) {
	nd2, ok := g[n2]
	if !ok {
		nd2 = &Node{Data: n2}
		g[n2] = nd2
	}
	if nd1, ok := g[n1]; !ok {
		nd1 = &Node{n1, []graph.Neighbor{{graph.Edge(ed), nd2}}}
		g[n1] = nd1
	} else {
		nd1.Nbs = append(nd1.Nbs, graph.Neighbor{graph.Edge(ed), nd2})
	}
}

func (n *Node) LinkFrom(prev graph.NeighborNode, ed graph.Edge) graph.NeighborNode {
	rn := &Node{Data: n} // create new node referring to receiver.
	if prev != nil {
		nb := graph.Neighbor{Nd: rn}
		if de, ok := ed.(graph.DistanceEdge); ok {
			nb.Ed = Edge(de.Distance()) // create edge if meaningful
		}
		pn := prev.(*Node)
		pn.Nbs = append(pn.Nbs, nb)
	}
	return rn
}
