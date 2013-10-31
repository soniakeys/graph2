// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

// Graph implements Dijkstra's shortest path algorithm.
//
// Dijkstra's algorithm is efficient for finding the shortest path between
// two nodes in a general directed or undirected graph.  The path length
// minimized is the sum of edge lengths in the path, which must be non-negative.
package graph

import "container/heap"

// Node methods return data on a node (or vertex.)
type Node interface {
	// Neighbors returns neighbors of this node; that is, nodes reachable
	// from this one by following a single edge.  For an undirected graph,
	// if node B is in the list returned by A.Neighbors() then A must be in
	// the list returned by B.Neighbors().  Further, the edge returned in
	// each case must be the same, at least by the result of Edge.Distance().
	//
	// For efficiency concerns, ideally an implentation will maintain a
	// neighbor list for each node as part of its graph representation and
	// simply return the list.  In this case implementations will ignore
	// the argument to Neighbors.
	//
	// The argument however, is useful in cases where the implementation
	// constructs the neighbor list on demand, which might be reasonable in
	// some cases.  The argument passed by DijkstraShortestPath is always
	// a zero length slice but may be the previous result truncated to zero
	// length.  Implementations can append to this slice, reusing the
	// existing capacity and minimizing garbage.
	Neighbors([]Neighbor) []Neighbor
	// D must return a pointer to a Dijkstra struct that is assoiated with
	// the node.  Each node must have its own Dijkstra struct.
	D() *Dijkstra
}

// Edge returns data on an edge.  The only data needed is the edge distance
// (or length, or weight.)  This distance must be non-negative.
type Edge interface {
	Distance() float64
}

// Neighbor describes an edge leading from a node and the node at the other
// end of the edge.  It thus represents a directed link.  In an undirected
// graph, a neighbor of one node always has the first node as one of its
// neighbors.
type Neighbor struct {
	Edge
	Node
}

// Dijkstra holds data needed internally by the package implementation of
// Dijkstra's algorithm.  Implementations of Node need to create one of
// these structs for each node and maintain the struct associated with the
// node.
type Dijkstra struct {
	tent     float64 // tentative distance
	prevNode Node    // path back to start
	prevEdge Edge    // edge from prevNode to the node of this struct
	n        int     // number of nodes in path
	done     bool    // true when tent and prev represent shortest path
	rx       int     // heap.Remove index
}

// Reset prepares a Dijkstra struct for a shortest path search.  It should
// be called in the implementation of Graph.ResetDijkstra.
func (d *Dijkstra) Reset() {
	d.n = 0
	d.done = false
}

// Graph represents all nodes in a graph.  The only method needed,
// ResetDijkstra, is called by DijkstraShortestPath as initialization.
// The implementation of ResetDijkstra must iterate over all nodes in the
// graph and call Reset() on the Dijkstra struct for each node.  There are
// no constraints on graph representation except that this function should
// have an efficient way of iterating over all nodes.
type Graph interface {
	ResetDijkstra()
}

// DijkstraShortestPath finds the shortest path between two nodes.
//
// It finds the shortest path between two nodes in a general directed or
// undirected graph.  The path length minimized is the sum of edge lengths
// in the path, which must be non-negative.
//
// Graph, Node, and Edge must be implemented as described in this package
// documentation. Argument g must be a properly connected graph, start and
// end must be nodes in g.  The found shortest path is returned as a Neighbor
// slice.  The first element of this slice will be the start node.  (The edge
// member will be nil, as there is no edge that needs be identified going to
// the start node.)  Remaining elements give the found path of edges and nodes.
// Also returned is the total path length.  If the end node cannot be reached
// from the start node, the returned neighbor list will be nil.
func DijkstraShortestPath(g Graph, start, end Node) ([]Neighbor, float64) {
	// WP steps 1 and 2.
	// WP references are to the algorithm description on Wikepedia,
	// http://en.wikipedia.org/wiki/Dijkstra%27s_algorithm#Algorithm
	g.ResetDijkstra()
	current := start
	cd := current.D()
	cd.n = 1
	cd.tent = 0
	var unvis ndList // heap
	var nb []Neighbor
	for {
		// WP step 3: update tentative distances to neighbors
		for _, nb := range current.Neighbors(nb[:0]) {
			if nd := nb.D(); !nd.done {
				d := cd.tent + nb.Distance()
				if nd.n > 0 {
					if nd.tent <= d {
						continue
					}
					heap.Remove(&unvis, nd.rx)
				}
				nd.tent = d
				nd.prevNode = current
				nd.prevEdge = nb.Edge
				nd.n = cd.n + 1
				heap.Push(&unvis, nb.Node)
			}
		}
		// WP step 4: mark current node visited, record path and distance
		cd.done = true
		if current == end {
			// WP step 5 (case of end node reached)
			// record path and distance for return value
			distance := cd.tent
			// recover path by tracing prev links
			i := cd.n
			path := make([]Neighbor, i)
			for n := current; n != nil; {
				i--
				d := n.D()
				path[i] = Neighbor{d.prevEdge, n}
				n = d.prevNode
			}
			return path, distance
		}
		if len(unvis) == 0 {
			break // WP step 5 (case of no more reachable nodes)
		}
		// WP step 6: new current is node with smallest tentative distance
		current = heap.Pop(&unvis).(Node)
		cd = current.D()
	}
	return nil, 0
}

// ndList implements container/heap
type ndList []Node

func (n ndList) Len() int           { return len(n) }
func (n ndList) Less(i, j int) bool { return n[i].D().tent < n[j].D().tent }
func (n ndList) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
	n[i].D().rx = i
	n[j].D().rx = j
}
func (n *ndList) Push(x interface{}) {
	nd := x.(Node)
	nd.D().rx = len(*n)
	*n = append(*n, nd)
}
func (n *ndList) Pop() interface{} {
	s := *n
	if len(s) == 0 {
		return nil
	}
	last := len(s) - 1
	r := s[last]
	*n = s[:last]
	return r
}
