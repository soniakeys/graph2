// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package adj_test

import (
	"fmt"

	"github.com/soniakeys/graph"
	"github.com/soniakeys/graph/adj"
	"github.com/soniakeys/graph/search"
)

func ExampleDigraph_dijkstraShortestPath() {
	g := adj.Digraph{}
	g.Link("a", "b", adj.Weighted(7))
	g.Link("a", "c", adj.Weighted(9))
	g.Link("a", "f", adj.Weighted(14))
	g.Link("b", "c", adj.Weighted(10))
	g.Link("b", "d", adj.Weighted(15))
	g.Link("c", "d", adj.Weighted(11))
	g.Link("c", "f", adj.Weighted(2))
	g.Link("d", "e", adj.Weighted(6))
	g.Link("e", "f", adj.Weighted(9))
	// echo initial conditions
	fmt.Println("Directed graph with", len(g), "nodes")
	// run Dijkstra's shortest path algorithm
	path, l := search.DijkstraShortestPath(g["a"], g["e"])
	fmt.Println(`Shortest path from node "a" to node "e":`, path)
	fmt.Println("Path length:", l)
	// Output:
	// Directed graph with 6 nodes
	// Shortest path from node "a" to node "e": [{<nil> a} {9 c} {11 d} {6 e}]
	// Path length: 26
}

func ExampleGraph_dijkstra() {
	g := adj.NewGraph()
	g.Link("a", "b", adj.Weighted(7))
	g.Link("a", "c", adj.Weighted(9))
	g.Link("a", "f", adj.Weighted(14))
	g.Link("b", "c", adj.Weighted(10))
	g.Link("b", "d", adj.Weighted(15))
	g.Link("c", "d", adj.Weighted(11))
	g.Link("c", "f", adj.Weighted(2))
	g.Link("d", "e", adj.Weighted(6))
	g.Link("e", "f", adj.Weighted(9))
	// echo initial conditions
	fmt.Println("Undirected graph with", len(g.Nodes), "nodes")
	// run Dijkstra's shortest path algorithm
	path, l := search.DijkstraShortestPath(g.Nodes["a"], g.Nodes["e"])
	fmt.Println("Shortest path:", path)
	fmt.Println("Path length:", l)
	// Output:
	// Undirected graph with 6 nodes
	// Shortest path: [{<nil> a} {9 c} {2 f} {9 e}]
	// Path length: 20
}

func ExampleDigraph_dijkstraAllPaths() {
	g := adj.Digraph{}
	g.Link("a", "b", adj.Weighted(7))
	g.Link("a", "c", adj.Weighted(9))
	g.Link("a", "f", adj.Weighted(14))
	g.Link("b", "c", adj.Weighted(10))
	g.Link("b", "d", adj.Weighted(15))
	g.Link("c", "d", adj.Weighted(11))
	g.Link("c", "f", adj.Weighted(2))
	g.Link("d", "e", adj.Weighted(6))
	g.Link("e", "f", adj.Weighted(9))
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
	pp("", search.DijkstraAllPaths(g["a"]))
	// Output:
	// a
	// a b
	// a c
	// a c f
	// a c d
	// a c d e
}
