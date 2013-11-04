// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

import (
	"fmt"

	"github.com/soniakeys/graph"
)

// node represents a node in a directed graph.  It represents directed edges
// from the node with the handy DijkstraNeighbor type from the graph package.
type node struct {
	nbs  []graph.DijkstraNeighbor // directed edges as DijkstraNeighbors
	name string                   // example application specific data
}

// edge is a simple number representing an edge length/distance/weight.
type edge float64

// node implements graph.DijkstraNode, also fmt.Stringer
func (n *node) Neighbors([]graph.DijkstraNeighbor) []graph.DijkstraNeighbor {
	return n.nbs
}
func (n *node) String() string { return n.name }

// edge implements graph.DijkstraEdge
func (e edge) Distance() float64 { return float64(e) }

// edgeData struct for simple specification of example data
type edgeData struct {
	v1, v2 string
	l      float64
}

// example data
var (
	exampleEdges = []edgeData{
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
	exampleStart = "a"
	exampleEnd   = "e"
)

// linkGraph constructs a linked representation of example data.
func linkGraph(g []edgeData, start, end string) (allNodes int, startNode, endNode *node) {
	all := map[string]*node{}
	// one pass over data to collect nodes
	for _, e := range g {
		if all[e.v1] == nil {
			all[e.v1] = &node{name: e.v1}
		}
		if all[e.v2] == nil {
			all[e.v2] = &node{name: e.v2}
		}
	}
	// second pass to link neighbors
	for _, ge := range g {
		n1 := all[ge.v1]
		n1.nbs = append(n1.nbs, graph.DijkstraNeighbor{edge(ge.l), all[ge.v2]})
	}
	return len(all), all[start], all[end]
}

func ExampleDijkstraShortestPath_directed() {
	// construct linked representation of example data
	allNodes, startNode, endNode :=
		linkGraph(exampleEdges, exampleStart, exampleEnd)
	// echo initial conditions
	fmt.Printf("Directed graph with %d nodes, %d edges\n",
		allNodes, len(exampleEdges))
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
