// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package search

import (
	"container/heap"
	"math"

	"github.com/soniakeys/graph2"
)

// DijkstraShortestPath finds a shortest path between two nodes.
//
// It finds a shortest path between two nodes in a general directed or
// undirected graph2.  The path length minimized is the sum of edge weights.
//
// Arguments start and end must implement graph2.AdjNode.  Edges connecting
// nodes must implement graph2.Weighted.  Weights must be non-negative and
// must not be an Inf or NaN.
//
// The found shortest path is returned as a graph2.Half slice.  The first
// element of this slice will be the start node.  (The edge member will be nil,
// as there is no edge that needs to be identified going to the start node.)
// Remaining elements give the found path of edges and nodes.
// Also returned is the total path length.  If the end node cannot be reached
// from the start node, the returned Half list will be nil and the path
// length +Inf.
func DijkstraShortestPath(start, end graph2.HalfNode) ([]graph2.Half, float64) {
	_, path, dist := djk(start, end, false)
	return path, dist
}

// DijkstraAllPaths finds the shortest paths from the start node to all other
// nodes in a graph, where path length is the sum of edge weights.
//
// Adjacency relationships between nodes can represent a general directed or
// undirected graph2.
//
// Edges connecting nodes must implement graph2.Weighted.
// Weights must be non-negative and must not be an Inf or NaN.
//
// The result map has a key for each node reachable from the start node.
// The element value of each key is a half edge representing the previous
// node along the shortest path.  The start node is included in the result,
// with a zero value element.
func DijkstraAllPaths(start graph2.HalfNode) map[graph2.HalfNode]graph2.FromHalf {
	tree, _, _ := djk(start, nil, true)
	return tree
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
	// status/index of tentPath in pool that backs heap
	tx int
}

// tentPath holds additional data for a node in the "tentative set".
type tentPath struct {
	dist float64 // tentative path distance
	n    int     // number of nodes in path
	rx   int     // heap.Remove index
	nd   graph2.HalfNode
}

type tentHeap struct {
	pool []tentPath
	heap []int // values are indexes into pool
	free []int // values are indexes into pool
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

func djk(start, end graph2.HalfNode, all bool) (map[graph2.HalfNode]graph2.FromHalf, []graph2.Half, float64) {
	if start == nil {
		return nil, nil, math.Inf(1)
	}
	current := start
	cd := dijkstra{tx: -1} // mark start done.  it skips the heap.
	d := map[graph2.HalfNode]dijkstra{start: cd}
	prev := map[graph2.HalfNode]graph2.FromHalf{start: graph2.FromHalf{}}
	ct := tentPath{n: 1} // path length 1 for start node
	h := &tentHeap{
		pool: make([]tentPath, 1)} // zero element unused
	for {
		if current == end { // single path search complete
			distance := ct.dist
			// recover path by tracing prev links
			i := ct.n
			path := make([]graph2.Half, i)
			for i > 0 {
				i--
				from := prev[current]
				path[i].Ed = from.Ed
				path[i].To = current
				current = from.From
			}
			return nil, path, distance // success
		}
		current.VisitAdjHalfs(func(a graph2.Half) {
			nd := d[a.To]
			if nd.tx < 0 {
				return // skip nodes already done
			}
			dist := ct.dist + a.Ed.(graph2.Weighted).Weight()
			if nd.tx > 0 { // node already in tentative set
				nt := &h.pool[nd.tx]
				if dist >= nt.dist {
					return // it's no help
				}
				// the path through current to this node is shorter than some
				// other path to this node.  record new path data and reheap.
				nt.dist = dist
				nt.n = ct.n + 1
				prev[a.To] = graph2.FromHalf{current, a.Ed}
				d[a.To] = nd
				heap.Fix(h, nt.rx)
			} else { // nd.tx was zero. this is the first visit to this node.
				// first find a place for tentPath data
				if len(h.free) == 0 {
					// nothing on the free list, extend the pool.
					nd.tx = len(h.pool)
					h.pool = append(h.pool, tentPath{
						nd:   a.To,
						dist: dist,
						n:    ct.n + 1})
				} else { // reuse
					last := len(h.free) - 1
					nd.tx = h.free[last]
					h.free = h.free[:last]
					h.pool[nd.tx] = tentPath{
						nd:   a.To,
						dist: dist,
						n:    ct.n + 1}
				}
				// push path data to heap
				prev[a.To] = graph2.FromHalf{current, a.Ed}
				d[a.To] = nd
				heap.Push(h, nd.tx)
			}
		})
		if len(h.heap) == 0 {
			//			return stRoot, nil, math.Inf(1)
			return prev, nil, math.Inf(1)
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
