// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package search_test

import (
	"math"
	"math/rand"
	"sort"
	"testing"

	"github.com/soniakeys/graph"
	"github.com/soniakeys/graph/search"
)

// minimal node type with no adjacency representation.
type djNode0 struct{}

func (n *djNode0) VisitAdjHalfs(graph.AdjHalfVisitor) {
}

func TestDijkstraDirected(t *testing.T) {
	// search from an isolated node to something other than itself
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

func (n *stNode) VisitAdjHalfs(v graph.AdjHalfVisitor) {
	for _, a := range n.nbs {
		v(graph.Half{a, a.to})
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
