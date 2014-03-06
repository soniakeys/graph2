// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package adj_test

import (
	"fmt"
	"sort"

	"github.com/soniakeys/graph"
	"github.com/soniakeys/graph/adj"
	"github.com/soniakeys/graph/search"
)

func ExampleDigraph_depthFirst() {
	g := adj.Digraph{}
	g.Link(5, 6, nil)
	g.Link(5, 7, nil)
	g.Link(5, 9, nil)
	g.Link(6, 7, nil)
	g.Link(7, 5, nil)
	g.Link(7, 7, nil)
	g.Link(7, 8, nil)
	v := func(n graph.Node, level int) bool {
		num := n.(*adj.Node).Data.(int)
		if num == 9 {
			return false
		}
		fmt.Println(num, "   ", level)
		return true
	}
	fmt.Println("Node  Level")
	fmt.Println(search.DepthFirst(g[5], v))
	// Output:
	// Node  Level
	// 5     0
	// 6     1
	// 7     2
	// 8     3
	// false
}

func ExampleDigraph_breadthFirst() {
	g := adj.Digraph{}
	g.Link(5, 6, nil)
	g.Link(5, 7, nil)
	g.Link(5, 9, nil)
	g.Link(6, 7, nil)
	g.Link(7, 5, nil)
	g.Link(7, 7, nil)
	g.Link(7, 8, nil)
	fmt.Println("Node  Level")
	v := func(n graph.Node, level int) bool {
		num := n.(*adj.Node).Data.(int)
		if num == 8 {
			return false
		}
		fmt.Println(num, "   ", level)
		return true
	}
	_, ok := search.BreadthFirst1(g[5], v)
	fmt.Println(ok)
	// Output:
	// Node  Level
	// 5     0
	// 6     1
	// 7     1
	// 9     1
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
	// run Dijkstra's algorithm to find all shortest paths
	from := search.DijkstraAllPaths(g["a"])
	// format output by walking each node of the result back to start
	as := make([]string, len(from))
	i := 0
	for nd, fh := range from {
		s := fmt.Sprint(nd)
		for fh.From != nil {
			s = fmt.Sprintf("%s %g %s", fh.From, fh.Ed, s)
			fh = from[fh.From]
		}
		as[i] = s
		i++
	}
	// sort for test repeatability
	sort.Strings(as)
	for _, s := range as {
		fmt.Println(s)
	}
	// Output:
	// a
	// a 7 b
	// a 9 c
	// a 9 c 11 d
	// a 9 c 11 d 6 e
	// a 9 c 2 f
}
