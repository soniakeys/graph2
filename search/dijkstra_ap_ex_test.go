// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package search_test

import (
	"fmt"
	"sort"

	"github.com/soniakeys/graph"
	"github.com/soniakeys/graph/search"
)

// DijkstraAllPaths requires a node type that implements graph.HalfNode
// and an edge type that implements graph.Weighted.  Our two types:

type (
	dapNode struct {
		name string       // node name
		nbs  []graph.Half // "neighbors," adjacent arcs and nodes
	}
	dapArc float64
)

// One method implements graph.HalfNode.
func (n dapNode) VisitAdjHalfs(v graph.AdjHalfVisitor) {
	for _, a := range n.nbs {
		v(a)
	}
}

// One method implements graph.Weighted.
func (a dapArc) Weight() float64 {
	return float64(a)
}

// Another method on dapNode implements fmt.Stringer
func (n dapNode) String() string {
	return n.name
}

// One more method on dapNode to make graph construction easy.
func (n *dapNode) link(n2 *dapNode, weight int) {
	n.nbs = append(n.nbs, graph.Half{dapArc(weight), n2})
}

func ExampleDijkstraAllPaths() {
	a := &dapNode{name: "a"}
	b := &dapNode{name: "b"}
	c := &dapNode{name: "c"}
	d := &dapNode{name: "d"}
	e := &dapNode{name: "e"}
	f := &dapNode{name: "f"}
	a.link(b, 7)
	a.link(c, 9)
	a.link(f, 14)
	b.link(c, 10)
	b.link(d, 15)
	c.link(d, 11)
	c.link(f, 2)
	d.link(e, 6)
	e.link(f, 9)
	// run Dijkstra's algorithm to find all shortest paths
	from := search.DijkstraAllPaths(a)
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
