// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

import (
	"math"
	"math/rand"
	"sort"
	"testing"

	"github.com/soniakeys/graph"
)

type stNode struct {
	name string
	x, y float64
	nbs  []stEdge
}

type stEdge struct {
	length float64
	to     *stNode
}

func (n *stNode) DistanceNeighbors(nbs []graph.DistanceNeighbor) []graph.DistanceNeighbor {
	for _, e := range n.nbs {
		nbs = append(nbs, graph.DistanceNeighbor{e, e.to})
	}
	return nbs
}

func (n *stNode) String() string { return n.name }

func (e stEdge) Distance() float64 { return float64(e.length) }

type xyList []stNode

func (l xyList) Len() int           { return len(l) }
func (l xyList) Less(i, j int) bool { return l[i].name < l[j].name }
func (l xyList) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }

// generate a random graph
func r(nNodes, nEdges int) (start, end *stNode) {
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
	for i := 0; i < nEdges; {
		n1 := &nodes[s.Intn(nNodes)]
		n2 := &nodes[s.Intn(nNodes)]
		dist := math.Hypot(n2.x-n1.x, n2.y-n1.y)
		if dist > s.Float64()*math.Sqrt2 {
			continue
		}
		n1.nbs = append(n1.nbs, stEdge{dist, n2})
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
		graph.DijkstraShortestPath(start, end)
	}
}

func Benchmark1e3(b *testing.B) {
	// 1000 nodes
	start, end := r(1000, 3000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		graph.DijkstraShortestPath(start, end)
	}
}

func Benchmark1e4(b *testing.B) {
	// 10k nodes
	start, end := r(1e4, 5e4)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		graph.DijkstraShortestPath(start, end)
	}
}

func Benchmark1e5(b *testing.B) {
	// 100k nodes
	start, end := r(1e5, 1e6)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		graph.DijkstraShortestPath(start, end)
	}
}
