package search_test

import (
	"fmt"

	"github.com/soniakeys/graph"
	"github.com/soniakeys/graph/search"
)

type bfsNode struct {
	num int
	nbs []graph.Node
}

func (n bfsNode) VisitOk(v graph.NodeOkVisitor) bool {
	for _, nb := range n.nbs {
		if !v(nb) {
			return false
		}
	}
	return true
}

func (n bfsNode) NumAdj() int {
	return len(n.nbs)
}

func ExampleBreadthFirstSimple() {
	n0 := &bfsNode{num: 0}
	n1 := &bfsNode{num: 1}
	n2 := &bfsNode{num: 2}
	n3 := &bfsNode{num: 3}
	n4 := &bfsNode{num: 4}
	n0.nbs = []graph.Node{n1, n2, n4}
	n1.nbs = []graph.Node{n2}
	n2.nbs = []graph.Node{n0, n2, n3}
	v := func(n graph.Node) bool {
		num := n.(*bfsNode).num
		if num == 3 {
			return false
		}
		fmt.Println(num)
		return true
	}
	fmt.Println(search.BreadthFirstSimple(n0, v))
	// Output:
	// 0
	// 1
	// 2
	// 4
	// false
}
