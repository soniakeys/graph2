// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package adj

import (
	"fmt"

	"github.com/soniakeys/graph"
)

// The example contains an example graph representation and example data.
// The representation implements graph.EstimateNode and graph.DistanceEdge.

// estNode represents a node that can return a heuristic distance estimate
// to another node.  It represents directed edges from the node with the handy
// EstimateNeighbor type from the graph package.
type Node struct {
	Data interface{}      // application specific
	Nbs  []graph.Neighbor // directed edges
}

// Node implements graph.EstimateNode, also fmt.Stringer
func (n *Node) Visit(v graph.NeighborVisitor) {
	for _, nb := range n.Nbs {
		v(nb)
	}
}
func (n *Node) Estimate(e graph.EstimateNode) float64 {
	return n.Data.(graph.Estimator).Estimate(e)
}
func (n *Node) String() string { return fmt.Sprint(n.Data) }

// Edge implements graph.DistanceEdge
type Edge float64

func (e Edge) Distance() float64 { return float64(e) }

type Graph map[interface{}]*Node

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
