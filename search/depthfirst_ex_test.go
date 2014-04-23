package search_test

import (
	"fmt"

	"github.com/soniakeys/graph2"
	"github.com/soniakeys/graph2/search"
)

// dfNode implements graph2.Node.
type dfNode struct {
	num int
	nbs []graph2.Node
}

// VisitAdjNodes is the only method needed to satisfy the interface.
func (n dfNode) VisitAdjNodes(v graph2.AdjNodeVisitor) bool {
	for _, nb := range n.nbs {
		if !v(nb) {
			return false
		}
	}
	return true
}

func ExampleDepthFirst() {
	n5 := &dfNode{num: 5}
	n6 := &dfNode{num: 6}
	n7 := &dfNode{num: 7}
	n8 := &dfNode{num: 8}
	n9 := &dfNode{num: 9}
	n5.nbs = []graph2.Node{n6, n7, n9}
	n6.nbs = []graph2.Node{n7}
	n7.nbs = []graph2.Node{n5, n7, n8}
	fmt.Println("Node  Level")
	v := func(n graph2.Node, level int) bool {
		num := n.(*dfNode).num
		if num == 9 {
			return false
		}
		fmt.Println(num, "    ", level)
		return true
	}
	fmt.Println(search.DepthFirst(n5, v))
	// Output:
	// Node  Level
	// 5      0
	// 6      1
	// 7      2
	// 8      3
	// false
}
