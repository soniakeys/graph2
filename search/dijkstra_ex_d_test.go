// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package search_test

import (
	"fmt"

	"github.com/soniakeys/graph/adj"
	"github.com/soniakeys/graph/search"
)

func ExampleDijkstraShortestPath_directed() {
	g := adj.Graph{}
	g.Link("a", "b", adj.Edge(7))
	g.Link("a", "c", adj.Edge(9))
	g.Link("a", "f", adj.Edge(14))
	g.Link("b", "c", adj.Edge(10))
	g.Link("b", "d", adj.Edge(15))
	g.Link("c", "d", adj.Edge(11))
	g.Link("c", "f", adj.Edge(2))
	g.Link("d", "e", adj.Edge(6))
	g.Link("e", "f", adj.Edge(9))
	// echo initial conditions
	fmt.Println("Directed graph with", len(g), "nodes")
	// run Dijkstra's shortest path algorithm
	ndPath, edPath, l := search.DijkstraShortestPath(g["a"], g["e"])
	if ndPath == nil {
		fmt.Println(`No path from node "a" to node "e"`)
		return
	}
	fmt.Println(`Shortest path from node "a" to node "e":`, ndPath)
	fmt.Println("Edges linking path nodes:", edPath)
	fmt.Println("Path length:", l)
	// Output:
	// Directed graph with 6 nodes
	// Shortest path from node "a" to node "e": [a c d e]
	// Edges linking path nodes: [9 11 6]
	// Path length: 26
}
