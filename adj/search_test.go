// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package adj_test

import (
	"fmt"

	"github.com/soniakeys/graph"
	"github.com/soniakeys/graph/adj"
	"github.com/soniakeys/graph/search"
)

func ExampleDigraph_depthFirst() {
	g := adj.Digraph{}
	g.Link(0, 1, nil)
	g.Link(0, 2, nil)
	g.Link(0, 4, nil)
	g.Link(1, 2, nil)
	g.Link(2, 0, nil)
	g.Link(2, 2, nil)
	g.Link(2, 3, nil)
	v := func(n graph.Node) bool {
		num := n.(*adj.Node).Data.(int)
		if num == 4 {
			return false
		}
		fmt.Println(num)
		return true
	}
	fmt.Println(search.DepthFirst(g[0], v))
	// Output:
	// 0
	// 1
	// 2
	// 3
	// false
}

func ExampleDigraph_breadthFirst() {
	g := adj.Digraph{}
	g.Link(0, 1, nil)
	g.Link(0, 2, nil)
	g.Link(0, 4, nil)
	g.Link(1, 2, nil)
	g.Link(2, 0, nil)
	g.Link(2, 2, nil)
	g.Link(2, 3, nil)
	v := func(n graph.Node) bool {
		num := n.(*adj.Node).Data.(int)
		if num == 3 {
			return false
		}
		fmt.Println(num)
		return true
	}
	fmt.Println(search.BreadthFirst(g[0], v))
	// Output:
	// 0
	// 1
	// 2
	// 4
	// false
}

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
		n.VisitAdj(func(nb graph.Half) {
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
