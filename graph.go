// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

// Graph implements Dijkstra's shortest path algorithm.
//
// Dijkstra's algorithm is efficient for finding the shortest path between
// two nodes in a general directed or undirected graph.  The path length
// minimized is the sum of edge lengths in the path, which must be non-negative.
package graph

import "container/heap"

// DijkstraNode method returns data on a node (or vertex) as required by
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

// dijkstra holds data per node that is needed by the algorithm.  The
// data is associated with nodes with a map in the search struct (below.)
// Data is not added to the map until a node is "visited" or reached by
// the search.  Field tx represents the status of a node which may be
// unvisited, tentative, or done.  tx = 0 represents unvisited so that
// the zero value of the dijkstra type--returned when a node is not in
// the map yet--will reflect this status.  A node first reached from another
// node is then moved to the "tentative set," a set of nodes maintained
// as a min heap.  Additional data (tentPath) is needed for nodes in the
// tentative set.  tx is then used for a 1-based index to this additional data.
// when a node is removed from the heap, tx is set to -1, indicating "done,"
// and prevNode and prevEdge are updated with will values representing the
// shortest path from the start node.
type dijkstra struct {
	tx       int          // status/index of tentPath in pool that backs heap
	prevNode DijkstraNode // path back to start
	prevEdge DijkstraEdge // edge from prevNode to the node of this struct
}

// tentPath holds additional data for a node in the "tentative set".
type tentPath struct {
	dist float64 // tentative path distance
	n    int     // number of nodes in path
	rx   int     // heap.Remove index
	nd   DijkstraNode
}

type tentHeap struct {
	pool []tentPath
	heap []int // values are indexes into pool
	free []int // values are indexes into pool
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
	current := start
	cd := dijkstra{tx: -1} // mark start done.  it skips the heap.
	d := map[DijkstraNode]dijkstra{current: cd}
	ct := tentPath{n: 1} // path length 1 for start node
	h := &tentHeap{
		pool: make([]tentPath, 1)} // zero element unused
	var nbs []DijkstraNeighbor // recycled slice
	for {
		nbs = current.Neighbors(nbs[:0])
		for _, nb := range nbs {
			nd := d[nb.DijkstraNode]
			if nd.tx < 0 {
				continue // skip nodes already done
			}
			dist := ct.dist + nb.Distance()
			if nd.tx > 0 { // node already in tentative set
				nt := &h.pool[nd.tx]
				if dist >= nt.dist {
					continue // it's no help
				}
				// the path through current to this node is shorter than some
				// other path to this node.  record new path data and reheap.
				nt.dist = dist
				nt.n = ct.n + 1
				nd.prevNode = current
				nd.prevEdge = nb.DijkstraEdge
				d[nb.DijkstraNode] = nd
				heap.Fix(h, nt.rx)
			} else { // nd.tx was zero. this is the first visit to this node.
				// first find a place for tentPath data
				if len(h.free) == 0 {
					// nothing on the free list, extend the pool.
					nd.tx = len(h.pool)
					h.pool = append(h.pool, tentPath{
						nd:   nb.DijkstraNode,
						dist: dist,
						n:    ct.n + 1})
				} else { // reuse
					last := len(h.free) - 1
					nd.tx = h.free[last]
					h.free = h.free[:last]
					h.pool[nd.tx] = tentPath{
						nd:   nb.DijkstraNode,
						dist: dist,
						n:    ct.n + 1}
				}
				// push path data to heap
				nd.prevNode = current
				nd.prevEdge = nb.DijkstraEdge
				d[nb.DijkstraNode] = nd
				heap.Push(h, nd.tx)
			}
		}
		if current == end { // search complete
			distance := ct.dist
			// recover path by tracing prev links
			i := ct.n
			path := make([]DijkstraNeighbor, i)
			for n := current; n != nil; {
				i--
				nd := d[n]
				path[i] = DijkstraNeighbor{nd.prevEdge, n}
				n = nd.prevNode
			}
			return path, distance // success
		}
		if len(h.heap) == 0 {
			return nil, 0 // failure. no more reachable nodes
		}
		// new current is node with smallest tentative distance
		ctx := heap.Pop(h).(int)
		ct = h.pool[ctx]
		current = ct.nd
		cd = d[current]
		h.free = append(h.free, ctx) // recycle tentPath struct
		cd.tx = -1                   // done
		d[current] = cd              // store the -1
	}
}

// search implements container/heap
func (h tentHeap) Len() int { return len(h.heap) }
func (h tentHeap) Less(i, j int) bool {
	return h.pool[h.heap[i]].dist < h.pool[h.heap[j]].dist
}
func (h tentHeap) Swap(i, j int) {
	h.heap[i], h.heap[j] = h.heap[j], h.heap[i]
	h.pool[h.heap[i]].rx = i
	h.pool[h.heap[j]].rx = j
}
func (h *tentHeap) Push(x interface{}) {
	tx := x.(int)
	h.pool[tx].rx = len(h.heap)
	h.heap = append(h.heap, tx)
}
func (h *tentHeap) Pop() interface{} {
	last := len(h.heap) - 1
	tx := h.heap[last]
	h.heap = h.heap[:last]
	return tx
}
