package search_test

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/soniakeys/graph"
	"github.com/soniakeys/graph/search"
)

type bfuNode struct {
	num int
	nbs []graph.BF2Node
}

func (n bfuNode) VisitBF2In(v graph.BF2NeighborVisitor) bool {
	for _, nb := range n.nbs {
		switch v(nb) {
		case graph.BF2Stop:
			return false
		case graph.BF2Found:
			return true
		}
	}
	return true
}

func (n bfuNode) VisitBF2Out(v graph.BF2NeighborVisitor) bool {
	for _, nb := range n.nbs {
		if v(nb) == graph.BF2Stop {
			return false
		}
	}
	return true
}

func (n bfuNode) String() string {
	return fmt.Sprint(n.num)
}

func (n bfuNode) NumAdj() int {
	return len(n.nbs)
}

type bfuGraph struct {
	nds map[int]*bfuNode
	ned int
}

func (g bfuGraph) NumEdges() int { return g.ned }
func (g bfuGraph) Nodes() map[graph.BF2Node]struct{} {
	m := make(map[graph.BF2Node]struct{}, len(g.nds))
	for _, n := range g.nds {
		m[n] = struct{}{}
	}
	return m
}

func (g *bfuGraph) link(n1, n2 int) {
	nd1, ok := g.nds[n1]
	if !ok {
		nd1 = &bfuNode{num: n1}
		g.nds[n1] = nd1
	}
	nd2, ok := g.nds[n2]
	if !ok {
		nd2 = &bfuNode{num: n2}
		g.nds[n2] = nd2
	}
	nd1.nbs = append(nd1.nbs, nd2)
	nd2.nbs = append(nd2.nbs, nd1)
	g.ned++
}

func TestBF2_undirected(t *testing.T) {
	g := &bfuGraph{nds: map[int]*bfuNode{}}
	g.link(1, 3)
	g.link(3, 5)
	g.link(2, 5)
	g.link(4, 5)
	g.link(6, 5)
	g.link(9, 5)
	g.link(10, 5)
	g.link(11, 5)
	g.link(12, 5)
	g.link(12, 13)
	g.link(11, 14)
	g.link(11, 15)
	g.link(14, 15)
	g.link(26, 15)
	g.link(26, 5)
	g.link(26, 30)
	g.link(26, 29)
	g.link(26, 28)
	g.link(26, 27)
	g.link(26, 16)
	g.link(10, 16)
	g.link(17, 16)
	g.link(17, 10)
	g.link(17, 9)
	g.link(19, 9)
	g.link(19, 5)
	g.link(19, 8)
	g.link(7, 8)
	g.link(19, 18)
	g.link(19, 20)
	g.link(19, 22)
	g.link(19, 21)
	g.link(25, 21)
	g.link(24, 21)
	g.link(23, 21)
	g.link(24, 31)
	g.link(25, 32)
	g.link(33, 32)
	// test is to list paths up to two levels from node 17
	v := func(n graph.BF2Node, l int) (ok bool) {
		return l <= 2
	}
	start := g.nds[17]
	p, ok := search.BreadthFirst2(g, start, v)
	if ok { // expecing !ok after curtailing search at level 3
		t.Fatal(ok)
	}
	paths := make([]string, len(p))
	i := 0
	for n1, n0 := range p {
		s := fmt.Sprint(n1)
		for n1 != start {
			n1 = n0
			n0 = p[n0]
			s = fmt.Sprintf("%s %s", n1, s)
		}
		paths[i] = s
		i++
	}
	sort.Strings(paths)
	s := strings.Join(paths, "\n")
	want := `17
17 10
17 16
17 16 26
17 9
17 9 19
17 9 5`
	if s != want {
		t.Fatalf(`got
%s
want:
%s
`, s, want)
	}
}
