// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package search_test

import (
	"fmt"

	"github.com/soniakeys/graph"
	"github.com/soniakeys/graph/search"
)

// uNode represents a node in an undirected graph.
type uNode struct {
	eds  []*uEdge // list of neighboring edges
	name string   // example application specific data
}

// uEdge represents an undirected edge.
type uEdge struct {
	n1, n2 *uNode  // each edge connects two nodes
	dist   float64 // used to implement Distance method
}

func (e *uEdge) opp(n *uNode) *uNode {
	if n == e.n1 {
		return e.n2
	}
	return e.n1
}

// uNode implements graph.DistanceNode, also fmt.Stringer
// dxNode implements graph.NeighborNode, also fmt.Stringer
func (n *uNode) Visit(v graph.NeighborVisitor) {
	for _, ed := range n.eds {
		v(graph.Neighbor{ed, ed.opp(n)})
	}
}
func (n *uNode) String() string { return n.name }

// uEdge implements graph.DistanceEdge, also fmt.Stringer
func (e uEdge) String() string    { return fmt.Sprint(e.dist) }
func (e uEdge) Distance() float64 { return e.dist }

// uEdgeData struct for simple specification of example data
type uEdgeData struct {
	v1, v2 string
	l      float64
}

// example data
var (
	nd = []string{"a", "b", "c", "d", "e", "f"}
	ed = []uEdgeData{
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
	uStart = "a"
	uEnd   = "e"
)

// linkUGraph constructs a linked representation of example data.
func linkUGraph(nd []string, ed []uEdgeData, start, end string) (startUNode, endUNode *uNode) {
	all := make([]*uNode, len(nd))
	// construct nodes
	for i, n := range nd {
		all[i] = &uNode{name: n}
	}
	// link neighbors
	for _, ge := range ed {
		n1 := all[ge.v1[0]-'a']
		n2 := all[ge.v2[0]-'a']
		e := &uEdge{n1, n2, ge.l}
		n1.eds = append(n1.eds, e)
		n2.eds = append(n2.eds, e)
	}
	return all[start[0]-'a'], all[end[0]-'a']
}

func ExampleDijkstraShortestPath_undirected() {
	// construct linked representation of example data
	startNode, endNode := linkUGraph(nd, ed, uStart, uEnd)
	// echo initial conditions
	fmt.Printf("Undirected graph with %d nodes, %d edges\n", len(nd), len(ed))
	// run Dijkstra's shortest path algorithm
	ndPath, edPath, l := search.DijkstraShortestPath(startNode, endNode)
	if ndPath == nil {
		fmt.Println("No path from start node to end node")
		return
	}
	fmt.Println("Shortest path:", ndPath)
	fmt.Println("Edge distances:", edPath)
	fmt.Println("Path length:", l)
	// Output:
	// Undirected graph with 6 nodes, 9 edges
	// Shortest path: [a c f e]
	// Edge distances: [9 2 9]
	// Path length: 20
}
