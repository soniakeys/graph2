package graph_test

import (
	"fmt"

	"github.com/soniakeys/graph"
)

// uNode represents a node in an undirected graph.
type uNode struct {
	eds  []*uEdge        // list of neighboring edges
	d    *graph.Dijkstra // for D method of graph.Node interface
	name string          // example application specific data
}

// uEdge represents an undirected edge.
type uEdge struct {
	n1, n2 *uNode  // each edge connects two nodes
	dist   float64 // used to implement Distance method
	// more application specific data could go here
}

// uGraph represents an undirected graph.
type uGraph struct {
	nds      []*uNode // list of all nodes in the graph
	searched bool     // dirty bit directs reset
}

// uNode implements graph.Node, also fmt.Stringer
func (n *uNode) String() string     { return n.name }
func (n *uNode) D() *graph.Dijkstra { return n.d }
func (n *uNode) Neighbors(nbs []graph.Neighbor) []graph.Neighbor {
	for _, e := range n.eds {
		nb := graph.Neighbor{e, e.n1}
		if nb.Node == n {
			nb.Node = e.n2
		}
		nbs = append(nbs, nb)
	}
	return nbs
}

// uEdge implements graph.Edge, also fmt.Stringer
func (e uEdge) String() string    { return fmt.Sprint(e.dist) }
func (e uEdge) Distance() float64 { return e.dist }

// uGraph implements graph.Graph
func (g *uGraph) ResetDijkstra() {
	for _, n := range g.nds {
		if !g.searched {
			n.d = &graph.Dijkstra{}
		} else {
			n.D().Reset()
		}
	}
	g.searched = true
}

// uEdgeData struct for simple specification of example data
type uEdgeData struct {
	v1, v2 string
	l      float64
}

// example data
var (
	nd = []string{"a", "b", "c", "d", "e", "f"}
	ed = []uEdgeData{
		{"a", "b", 7},
		{"a", "c", 9},
		{"a", "f", 14},
		{"b", "c", 10},
		{"b", "d", 15},
		{"c", "d", 11},
		{"c", "f", 2},
		{"d", "e", 6},
		{"e", "f", 9},
	}
	uStart = "a"
	uEnd   = "e"
)

// linkUGraph constructs a linked representation of example data.
func linkUGraph(nd []string, ed []uEdgeData, start, end string) (allNodes *uGraph, startUNode, endUNode *uNode) {
	all := &uGraph{nds: make([]*uNode, len(nd))}
	// construct nodes
	for i, n := range nd {
		all.nds[i] = &uNode{name: n}
	}
	// link neighbors
	for _, ge := range ed {
		n1 := all.nds[ge.v1[0]-'a']
		n2 := all.nds[ge.v2[0]-'a']
		e := &uEdge{n1, n2, ge.l}
		n1.eds = append(n1.eds, e)
		n2.eds = append(n2.eds, e)
	}
	return all, all.nds[start[0]-'a'], all.nds[end[0]-'a']
}

func ExampleDijkstraShortestPath_undirected() {
	// construct linked representation of example data
	allNodes, startNode, endNode := linkUGraph(nd, ed, uStart, uEnd)
	// echo initial conditions
	fmt.Printf("Undirected graph with %d nodes, %d edges\n", len(nd), len(ed))
	// run Dijkstra's shortest path algorithm
	p, l := graph.DijkstraShortestPath(allNodes, startNode, endNode)
	if p == nil {
		fmt.Println("No path from start node to end node")
		return
	}
	fmt.Println("Shortest path:", p)
	fmt.Println("Path length:", l)
	// Output:
	// Undirected graph with 6 nodes, 9 edges
	// Shortest path: [{<nil> a} {9 c} {2 f} {9 e}]
	// Path length: 20
}
