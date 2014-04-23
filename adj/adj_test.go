package adj_test

import (
	"fmt"
	"sort"

	"github.com/soniakeys/graph2"
	"github.com/soniakeys/graph2/adj"
)

func ExampleNode_VisitAdjNodes() {
	g := adj.Digraph{}
	g.Link(0, 1, nil)
	g.Link(0, 2, nil)
	g.Link(0, 3, nil)
	v := func(n graph2.Node) bool {
		num := n.(*adj.Node).Data.(int)
		if num == 3 {
			return false
		}
		fmt.Println(num)
		return true
	}
	fmt.Println(g[0].VisitAdjNodes(v))
	// Output:
	// 1
	// 2
	// false
}

func ExampleWeighted_Weight() {
	// Example shows that adj.Weighted implements graph2.Weighted.
	var a graph2.Weighted
	a = adj.Weighted(4)
	fmt.Println(a.Weight())
	// Output:
	// 4
}

func ExampleDigraph_Link() {
	g := adj.Digraph{}

	// As a minimimal example, use ints for nodes and don't use arcs at all.
	g.Link(1, 2, nil)
	g.Link(2, 3, nil)
	g.Link(2, 1, nil)

	// Just for this example, buffer and sort output because maps are unordered.
	var output []string
	for id, nd := range g {
		// For each node, print the node.
		line := fmt.Sprintf("adjacent to node %v:", id)
		nd.VisitAdjHalfs(func(h graph2.Half) {
			// Print a list of neighbors on the same line.
			line += fmt.Sprintf(" %v", h.To)
		})
		output = append(output, line)
	}
	sort.Strings(output)
	for _, line := range output {
		fmt.Println(line)
	}
	// Output:
	// adjacent to node 1: 2
	// adjacent to node 2: 3 1
	// adjacent to node 3:
}
