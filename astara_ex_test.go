// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

import (
	"fmt"

	"github.com/soniakeys/graph"
)

// The example contains an example graph representation and example data.
// The representation implements graph.EstimateNode and graph.DistanceEdge.

// estNode represents a node that can return a heuristic distance estimate
// to another node.  It represents directed edges from the node with the handy
// EstimateNeighbor type from the graph package.
type estNode struct {
	nbs []graph.EstimateNeighbor // directed edges as DistanceNeighbors
	// hard coded distance estimate (left at 0.0 currently)
	hEnd float64
	name string // example application specific data
}

// estNode implements graph.EstimateNode, also fmt.Stringer
func (n *estNode) EstimateNeighbors([]graph.EstimateNeighbor) []graph.EstimateNeighbor {
	return n.nbs
}
func (n *estNode) Estimate(graph.EstimateNode) float64 { return n.hEnd }
func (n *estNode) String() string                      { return n.name }

// estEdge implements graph.DistanceEdge
type estEdge float64

func (e estEdge) Distance() float64 { return float64(e) }

var (
	estNodeData = []string{"a", "b", "c", "d", "e", "f"}
	estEdgeData = []struct {
		v1, v2 string
		l      int
	}{
		{"a", "b", 7},
		{"a", "c", 9},
		{"a", "f", 14},
		{"b", "c", 10},
		{"b", "d", 15},
		{"c", "d", 11},
		{"c", "f", 2},
		{"d", "e", 6},
		{"e", "f", 9},
	}
)

// linkEstGraph constructs a linked representation of example data.
func linkEstGraph() (startNode, endNode *estNode) {
	// create nodes
	all := map[string]*estNode{}
	for _, n := range estNodeData {
		all[n] = &estNode{name: n}
	}
	// link neighbors
	for _, ge := range estEdgeData {
		n1 := all[ge.v1]
		n1.nbs = append(n1.nbs,
			graph.EstimateNeighbor{estEdge(ge.l), all[ge.v2]})
	}
	return all["a"], all["e"]
}

func ExampleAStarA() {
	// construct linked representation of example data
	startNode, endNode := linkEstGraph()
	// echo initial conditions
	fmt.Printf("Directed graph with %d nodes, %d edges\n",
		len(estNodeData), len(estEdgeData))
	// run AStarA
	p, l := graph.AStarA(startNode, endNode)
	if p == nil {
		fmt.Println("No path from start node to end node")
		return
	}
	fmt.Println("Shortest path:", p)
	fmt.Println("Path length:", l)
	// Output:
	// Directed graph with 6 nodes, 9 edges
	// Shortest path: [{<nil> a} {9 c} {11 d} {6 e}]
	// Path length: 26
}
