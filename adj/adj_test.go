package adj_test

import (
	"fmt"
	"sort"

	"github.com/soniakeys/graph"
	"github.com/soniakeys/graph/adj"
)

func ExampleWeighted_Weight() {
	// Example shows that adj.Weighted implements graph.Weighted.
	var a graph.Weighted
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
		nd.Visit(func(h graph.Half) {
			// Print a list of neighbors on the same line.
			line += fmt.Sprintf(" %v", h.Nd)
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
