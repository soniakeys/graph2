package search_test

import (
	"fmt"

	"github.com/soniakeys/graph"
	"github.com/soniakeys/graph/search"
)

type dfNode struct {
	num int
	nbs []graph.Node
}

func (n dfNode) VisitAdjNodes(v graph.AdjNodeVisitor) bool {
	for _, nb := range n.nbs {
		if !v(nb) {
			return false
		}
	}
	return true
}

func ExampleDepthFirst() {
	n0 := &dfNode{num: 0}
	n1 := &dfNode{num: 1}
	n2 := &dfNode{num: 2}
	n3 := &dfNode{num: 3}
	n4 := &dfNode{num: 4}
	n0.nbs = []graph.Node{n1, n2, n4}
	n1.nbs = []graph.Node{n2}
	n2.nbs = []graph.Node{n0, n2, n3}
	v := func(n graph.Node, level int) bool {
		num := n.(*dfNode).num
		if num == 4 {
			return false
		}
		fmt.Println(num)
		return true
	}
	fmt.Println(search.DepthFirst(n0, v))
	// Output:
	// 0
	// 1
	// 2
	// 3
	// false
}
