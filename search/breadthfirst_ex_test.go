package search_test

import (
	"fmt"

	"github.com/soniakeys/graph2"
	"github.com/soniakeys/graph2/search"
)

type bfsNode struct {
	num int
	nbs []graph2.Node
}

func (n bfsNode) VisitAdjNodes(v graph2.AdjNodeVisitor) bool {
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

func ExampleBreadthFirst1() {
	n5 := &bfsNode{num: 5}
	n6 := &bfsNode{num: 6}
	n7 := &bfsNode{num: 7}
	n8 := &bfsNode{num: 8}
	n9 := &bfsNode{num: 9}
	n5.nbs = []graph2.Node{n6, n7, n9}
	n6.nbs = []graph2.Node{n7}
	n7.nbs = []graph2.Node{n5, n7, n8}
	fmt.Println("Node  Level")
	v := func(n graph2.Node, level int) bool {
		num := n.(*bfsNode).num
		if num == 8 {
			return false
		}
		fmt.Println(num, "   ", level)
		return true
	}
	_, ok := search.BreadthFirst1(n5, v)
	fmt.Println(ok)
	// Output:
	// Node  Level
	// 5     0
	// 6     1
	// 7     1
	// 9     1
	// false
}
