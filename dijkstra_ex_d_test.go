// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

import (
	"fmt"

	"github.com/soniakeys/graph"
)

// The example contains an example graph representation and example data.
// The representation implements graph.DistanceNode and graph.DistanceEdge.

// dxNode represents a node in a directed graph with non-negative edge
// distances.  It represents directed edges from the node with the handy
// DistanceNeighbor type from the graph package.
type dxNode struct {
	nbs  []graph.DistanceNeighbor // directed edges as DistanceNeighbors
	name string                   // example application specific data
}

// dxNode implements graph.DistanceNode, also fmt.Stringer
func (n *dxNode) DistanceNeighbors([]graph.DistanceNeighbor) []graph.DistanceNeighbor {
	return n.nbs
}
func (n *dxNode) String() string { return n.name }

// dxEdge implements graph.DistanceEdge
type dxEdge float64

func (e dxEdge) Distance() float64 { return float64(e) }

var (
	dxNodeData = []string{"a", "b", "c", "d", "e", "f"}
	dxEdgeData = []struct {
		v1, v2 string
		l      float64
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

// linkDxGraph constructs a linked representation of example data.
func linkDxGraph() (startNode, endNode *dxNode) {
	// create nodes
	all := map[string]*dxNode{}
	for _, n := range dxNodeData {
		all[n] = &dxNode{name: n}
	}
	// link neighbors
	for _, ge := range dxEdgeData {
		n1 := all[ge.v1]
		n1.nbs = append(n1.nbs,
			graph.DistanceNeighbor{dxEdge(ge.l), all[ge.v2]})
	}
	return all["a"], all["e"]
}

func ExampleDijkstraShortestPath_directed() {
	// construct linked representation of example data
	startNode, endNode := linkDxGraph()
	// echo initial conditions
	fmt.Printf("Directed graph with %d nodes, %d edges\n",
		len(dxNodeData), len(dxEdgeData))
	// run Dijkstra's shortest path algorithm
	p, l := graph.DijkstraShortestPath(startNode, endNode)
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
