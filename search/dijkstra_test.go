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

var s = rand.New(rand.NewSource(59))

// generate a random directed graph and end points to test
func r(nNodes, nArcs int, seed int64) (start, end *stNode) {
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
	// now generate random coordinates from repeatability seed
	s.Seed(seed)
	for i := range nodes {
		nodes[i].x = s.Float64()
		nodes[i].y = s.Float64()
	}
	// random start
	start = &nodes[s.Intn(nNodes)]
	// end is point at distance nearest target distance
	const target = .3
	nearest := 2.
	for i, n2 := range nodes {
		d := math.Abs(target - math.Hypot(n2.x-start.x, n2.y-start.y))
		if d < nearest {
			end = &nodes[i]
			nearest = d
		}
	}
generateArcs:
	for i := 0; i < nArcs; {
		n1 := s.Intn(nNodes)
		n2 := s.Intn(nNodes)
		nd1 := &nodes[n1]
		nd2 := &nodes[n2]
		dist := math.Hypot(nd2.x-nd1.x, nd2.y-nd1.y)
		if dist > s.ExpFloat64() {
			continue // favor near nodes
		}
		if n1 == n2 {
			continue // no self loops
		}
		for _, nb := range nd1.nbs {
			if nb.to == nd2 {
				continue generateArcs // no parallel arcs
			}
		}
		nd1.nbs = append(nd1.nbs, stArc{dist, nd2})
		i++
	}
	return
}

func Test100(t *testing.T) {
	start, end := r(100, 200, 62)
	t.Log(search.DijkstraShortestPath(start, end))
	t.Log(search.DijkstraShortestPath(start, end))
}

func Benchmark100(b *testing.B) {
	// 100 nodes
	start, end := r(100, 200, 62)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		search.DijkstraShortestPath(start, end)
	}
}

func Benchmark1e3(b *testing.B) {
	// 1000 nodes
	start, end := r(1000, 3000, 66)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		search.DijkstraShortestPath(start, end)
	}
}

func Benchmark1e4(b *testing.B) {
	// 10k nodes
	start, end := r(1e4, 5e4, 59)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		search.DijkstraShortestPath(start, end)
	}
}

func Benchmark1e5(b *testing.B) {
	// 100k nodes
	start, end := r(1e5, 1e6, 59)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		search.DijkstraShortestPath(start, end)
	}
}
