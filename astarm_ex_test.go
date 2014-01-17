// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

import (
	"fmt"

	"github.com/soniakeys/graph"
)

// The example contains an example graph representation and example data.
// The representation implements graph.EstimateNode and graph.DistanceEdge.

// monoNode represents a node that can return a heuristic distance estimate
// to another node.  It represents directed edges from the node with the handy
// EstimateNeighbor type from the graph package.
type monoNode struct {
	nbs  []graph.EstimateNeighbor // directed edges as EstimateNeighbors
	name string                   // example application specific data
	hEnd float64                  // heuristic distance estimate to end node
}

// monoNode implements graph.EstimateNode, also fmt.Stringer
func (n *monoNode) EstimateNeighbors([]graph.EstimateNeighbor) []graph.EstimateNeighbor {
	// validate monotonicity, for testing purposes
	for _, nb := range n.nbs {
		if n.Estimate(monoEnd) > nb.Distance()+nb.Estimate(monoEnd) {
			fmt.Printf(`non-monotonic:
	%s estimate = %f
	distance to %s = %f
	%s estimate = %f
	%f > %f + %f\n`,
				n, n.Estimate(monoEnd),
				nb.EstimateNode, nb.Distance(),
				nb.EstimateNode, nb.Estimate(monoEnd),
				n.Estimate(monoEnd), nb.Distance(), nb.Estimate(monoEnd))
		}
	}
	return n.nbs
}
func (n *monoNode) Estimate(graph.EstimateNode) float64 { return n.hEnd }
func (n *monoNode) String() string                      { return n.name }

// monoEdge implements graph.DistanceEdge
type monoEdge float64

func (e monoEdge) Distance() float64 { return float64(e) }

var (
	monoNodeData = []struct {
		name string
		hEnd int
	}{
		{"a", 6},
		{"b", 3},
		{"c", 4},
		{"d", 0},
/* WP example
		{"a", 19},
		{"b", 18},
		{"c", 10},
		{"d", 6},
		{"e", 0},
		{"f", 9},
*/
	}
	monoEdgeData = []struct {
		v1, v2 string
		l      int
	}{
		{"a", "b", 3},
		{"a", "c", 3},
		{"b", "d", 5},
		{"c", "b", 3},
		{"c", "d", 4},
/* WP example
		{"a", "b", 7},
		{"a", "c", 9},
		{"a", "f", 14},
		{"b", "c", 10},
		{"b", "d", 15},
		{"c", "d", 11},
		{"c", "f", 2},
		{"d", "e", 6},
		{"e", "f", 9},
*/
	}
)

// package variables for test inside EstimateNeighbors
var monoStart, monoEnd *monoNode

// linkMonoGraph constructs a linked representation of example data.
func linkMonoGraph() {
	// create nodes
	all := map[string]*monoNode{}
	for _, n := range monoNodeData {
		all[n.name] = &monoNode{name: n.name, hEnd: float64(n.hEnd)}
	}
	// link neighbors
	for _, ge := range monoEdgeData {
		n1 := all[ge.v1]
		n1.nbs = append(n1.nbs,
			graph.EstimateNeighbor{monoEdge(ge.l), all[ge.v2]})
	}
	monoStart = all["a"]
	monoEnd = all["d"]
}

func ExampleAStarM() {
	// construct linked representation of example data
	linkMonoGraph()
	// echo initial conditions
	fmt.Printf("Directed graph with %d nodes, %d edges\n",
		len(monoNodeData), len(monoEdgeData))
	// run AStarM
	p, l := graph.AStarM(monoStart, monoEnd)
	if p == nil {
		fmt.Println("No path from start node to end node")
		return
	}
	// verify admissability
	ap := 0.
	for i := len(p) - 1; ; {
		nd := p[i].EstimateNode
		if nd.Estimate(monoEnd) > ap {
			fmt.Printf(`inadmissable path:
	Estimate from %s was %f
	Actual path was %f
	%f > %f\n`,
				nd, nd.Estimate(monoEnd),
				ap,
				nd.Estimate(monoEnd), ap)
			return
		}
		if i == 0 {
			break
		}
		ap += p[i].Distance()
		i--
	}
	// good.
	fmt.Println("Shortest path:", p)
	fmt.Println("Path length:", l)
	// Output:
	// Directed graph with 4 nodes, 5 edges
	// Shortest path: [{<nil> a} {3 c} {4 d}]
	// Path length: 7
}
	// Shortest path: [{<nil> a} {9 c} {11 d} {6 e}]
	// Path length: 26
