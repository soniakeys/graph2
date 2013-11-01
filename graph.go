// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

// Graph implements Dijkstra's shortest path algorithm.
//
// Dijkstra's algorithm is efficient for finding the shortest path between
// two nodes in a general directed or undirected graph.  The path length
// minimized is the sum of edge lengths in the path, which must be non-negative.
package graph

import "container/heap"

// DijkstraNode methods return data on a node (or vertex) as required by
// DijkstraShortestPath.
type DijkstraNode interface {
	// Neighbors returns neighbors of this node; that is, nodes reachable
	// from this one by following a single edge.  For an undirected graph,
	// if node B is in the list returned by A.Neighbors() then A must be in
	// the list returned by B.Neighbors().  Further, the edge returned in
	// each case must be the same, at least by the result of Edge.Distance().
	//
	// For performance concerns, an implentation might maintain a neighbor
	// list for each node as part of its graph representation and simply
	// return the list.  In this case implementations will ignore the
	// argument to Neighbors.
	//
	// The argument however, is useful in cases where the implementation
	// constructs the neighbor list on demand.  The argument passed by
	// DijkstraShortestPath is always a zero length slice but may be the
	// previous result truncated to zero length.  Implementations can append
	// to this slice, reusing the existing capacity and minimizing garbage.
	Neighbors([]DijkstraNeighbor) []DijkstraNeighbor
}

// DijkstraEdge returns data on an edge.  The only data needed is the edge
// distance (or length, or weight.)  This distance must be non-negative.
type DijkstraEdge interface {
	Distance() float64
}

// DijkstraNeighbor describes an edge leading from a node and the node at
// the other end of the edge.  It thus represents a directed link.  In an
// undirected graph, a neighbor of one node always has the first node as
// one of its neighbors.
type DijkstraNeighbor struct {
	DijkstraEdge
	DijkstraNode
}

// dijkstra holds data for a node that has been visited.
type dijkstra struct {
	tp       *tentPath
	prevNode DijkstraNode // path back to start
	prevEdge DijkstraEdge // edge from prevNode to the node of this struct
	done     bool         // true when tent and prev represent shortest path
}

// tentPath holds additional data for a node in the "tentative set".
type tentPath struct {
	dist float64 // tentative path distance
	n    int     // number of nodes in path
	rx   int     // heap.Remove index
}

// search holds data needed to perform a single search.
type search struct {
	n []DijkstraNode // backs heap
	d map[DijkstraNode]dijkstra
}

// DijkstraShortestPath finds the shortest path between two nodes.
//
// It finds the shortest path between two nodes in a general directed or
// undirected graph.  The path length minimized is the sum of edge lengths
// in the path, which must be non-negative.
//
// DijkstraNode and DijkstraEdge must be implemented as described in this
// package documentation.  Arguments start and end must be nodes in a properly
// connected graph.  The found shortest path is returned as a DijkstraNeighbor
// slice.  The first element of this slice will be the start node.  (The edge
// member will be nil, as there is no edge that needs be identified going to
// the start node.)  Remaining elements give the found path of edges and nodes.
// Also returned is the total path length.  If the end node cannot be reached
// from the start node, the returned neighbor list will be nil.
func DijkstraShortestPath(start, end DijkstraNode) ([]DijkstraNeighbor, float64) {
	// WP steps 1 and 2.
	// WP references are to the algorithm description on Wikepedia,
	// http://en.wikipedia.org/wiki/Dijkstra%27s_algorithm#Algorithm
	current := start
	cd := dijkstra{tp: &tentPath{n: 1}}
	s := &search{d: map[DijkstraNode]dijkstra{current: cd}}
	var nbs []DijkstraNeighbor
	for {
		// WP step 3: update tentative distances to neighbors
		nbs = current.Neighbors(nbs[:0])
		for _, nb := range nbs {
			if nd := s.d[nb.DijkstraNode]; !nd.done {
				dist := cd.tp.dist + nb.Distance()
				tent := nd.tp != nil
				if tent && nd.tp.dist <= dist {
					continue
				}
				nd.prevNode = current
				nd.prevEdge = nb.DijkstraEdge
				if tent {
					nd.tp.dist = dist
					nd.tp.n = cd.tp.n + 1
					s.d[nb.DijkstraNode] = nd
					heap.Fix(s, nd.tp.rx)
				} else {
					nd.tp = &tentPath{
						dist: dist,
						n:    cd.tp.n + 1}
					s.d[nb.DijkstraNode] = nd
					heap.Push(s, nb.DijkstraNode)
				}
			}
		}
		// WP step 4: mark current node visited, record path and distance
		cd.done = true
		if current == end {
			// WP step 5 (case of end node reached)
			// record path and distance for return value
			distance := cd.tp.dist
			// recover path by tracing prev links
			i := cd.tp.n
			path := make([]DijkstraNeighbor, i)
			for n := current; n != nil; {
				i--
				nd := s.d[n]
				path[i] = DijkstraNeighbor{nd.prevEdge, n}
				n = nd.prevNode
			}
			return path, distance
		}
		if len(s.n) == 0 {
			break // WP step 5 (case of no more reachable nodes)
		}
		cd.tp = nil
		// WP step 6: new current is node with smallest tentative distance
		current = heap.Pop(s).(DijkstraNode)
		cd = s.d[current]
	}
	return nil, 0
}

// search implements container/heap
func (s search) Len() int { return len(s.n) }
func (s search) Less(i, j int) bool {
	return s.d[s.n[i]].tp.dist < s.d[s.n[j]].tp.dist
}
func (s search) Swap(i, j int) {
	s.n[i], s.n[j] = s.n[j], s.n[i]
	s.d[s.n[i]].tp.rx = i
	s.d[s.n[j]].tp.rx = j
}
func (s *search) Push(x interface{}) {
	nd := x.(DijkstraNode)
	s.d[nd].tp.rx = len(s.n)
	s.n = append(s.n, nd)
}
func (s *search) Pop() interface{} {
	if len(s.n) == 0 {
		return nil
	}
	last := len(s.n) - 1
	r := s.n[last]
	s.n = s.n[:last]
	return r
}
