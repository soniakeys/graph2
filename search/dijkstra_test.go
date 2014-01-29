// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package search_test

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"testing"

	"github.com/soniakeys/graph"
	"github.com/soniakeys/graph/adj"
	"github.com/soniakeys/graph/search"
)

func ExampleDijkstraShortestPath_directed() {
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

func ExampleDijkstraShortestPath_undirected() {
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

func ExampleDijkstraAllPaths() {
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
	var pp func(string, graph.NeighborNode)
	pp = func(s string, n graph.NeighborNode) {
		s += fmt.Sprint(n)
		fmt.Println(s)
		s += " "
		n.Visit(func(nb graph.Neighbor) {
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

// minimal node type with no neighbors.
type djNode0 struct{}

func (n *djNode0) Visit(graph.NeighborVisitor) {
}

func TestDijkstraDirected(t *testing.T) {
	// search from node with no neighbors to something other than itself
	// should return no path.
	p, l := search.DijkstraShortestPath(&djNode0{}, nil)
	if p != nil {
		t.Fatal("wrong node path.  expected nil, got", p)
	}
	if l != math.Inf(1) {
		t.Fatal("wrong path length.  expected +Inf, got", l)
	}
}

type stNode struct {
	name string
	x, y float64
	nbs  []stArc
}

type stArc struct {
	weight float64
	to     *stNode
}

func (n *stNode) Visit(v graph.NeighborVisitor) {
	for _, a := range n.nbs {
		v(graph.Neighbor{a, a.to})
	}
}

func (n *stNode) String() string { return n.name }

func (a stArc) Weight() float64 { return float64(a.weight) }

type xyList []stNode

func (l xyList) Len() int           { return len(l) }
func (l xyList) Less(i, j int) bool { return l[i].name < l[j].name }
func (l xyList) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }

// generate a random graph
func r(nNodes, nArcs int) (start, end *stNode) {
	s := rand.New(rand.NewSource(59))
	// generate unique node names
	nameMap := map[string]bool{}
	name := make([]byte, int(math.Log(float64(nNodes))/3)+1)
	for len(nameMap) < nNodes {
		for i := range name {
			name[i] = 'A' + byte(s.Intn(26))
		}
		nameMap[string(name)] = true
	}
	// sort for repeatability
	nodes := make(xyList, nNodes)
	i := 0
	for n := range nameMap {
		nodes[i].name = n
		i++
	}
	sort.Sort(nodes)
	// now assign random coordinates.
	for i := range nodes {
		nodes[i].x = s.Float64()
		nodes[i].y = s.Float64()
	}
	// generate edges.
	for i := 0; i < nArcs; {
		n1 := &nodes[s.Intn(nNodes)]
		n2 := &nodes[s.Intn(nNodes)]
		dist := math.Hypot(n2.x-n1.x, n2.y-n1.y)
		if dist > s.Float64()*math.Sqrt2 {
			continue
		}
		n1.nbs = append(n1.nbs, stArc{dist, n2})
		switch i {
		case 0:
			start = n1
		case 1:
			end = n2
		}
		i++
	}
	return start, end
}

func Benchmark100(b *testing.B) {
	// 100 nodes
	start, end := r(100, 200)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		search.DijkstraShortestPath(start, end)
	}
}

func Benchmark1e3(b *testing.B) {
	// 1000 nodes
	start, end := r(1000, 3000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		search.DijkstraShortestPath(start, end)
	}
}

func Benchmark1e4(b *testing.B) {
	// 10k nodes
	start, end := r(1e4, 5e4)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		search.DijkstraShortestPath(start, end)
	}
}

func Benchmark1e5(b *testing.B) {
	// 100k nodes
	start, end := r(1e5, 1e6)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		search.DijkstraShortestPath(start, end)
	}
}
