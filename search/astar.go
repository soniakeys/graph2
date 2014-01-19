// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package search

import (
	"container/heap"
	"math"

	"github.com/soniakeys/graph"
)

// AStarA finds a path between two nodes.
//
// AStarA implements both algorithm A and algorithm A*.  The difference in the
// two algorithms is strictly in the heuristic estimate returned by Estimate().
// Code for Estimate is not here but is provided by the caller.  If the
// caller provides an "admissable" heuristic estimate, then the algorithm is
// termed A*, otherwise it is algorithm A.  Admissable means the value returned
// by Estimate (an estimated distance from some node to a destination node)
// must be less than or equal to the actual shortest path distance from the
// node to the destination node.
//
// Like DijkstraShortestPath, AStarA with an admissable heuristic finds the
// shortest path between two nodes in a general directed or undirected graph
// with non-negative edge lengths.  Also edge lengths should not be Inf or NaN.
// AStarA generally runs faster than Dijkstra though, by using node the
// heuristic distance estimate.
//
// AStarA with an inadmissable heuristic becomes algorithm A.  Algorithm A
// will find a path, but it is not guaranteed to be the shortest path.
// The heuristic still guides the search however, so a nearly admissable
// heuristic is likely to find a very good path, if not the best.  Quality
// of the path returned degrades gracefully with the quality of the heuristic.
//
// Two interfaces, EstimateNode and DistanceEdge, must be implemented as
// described in this package documentation.  Arguments start and end must
// be nodes in a properly connected graph.
//
// The found path is returned as a graph.Neighbor slice.  The first
// element of this slice will be the start node.  (The edge member will be nil,
// as there is no edge that needs to be identified going to the start node.)
// Remaining elements give the found path of edges and nodes.
// Also returned is the total path length.  If the end node cannot be reached
// from the start node, the returned neighbor list will be nil and the path
// length +Inf.
func AStarA(start, end graph.EstimateNode) ([]graph.Neighbor, float64) {
	// start node is reached initially
	p := &rNode{
		nd: start,
		f:  start.Estimate(end),
		n:  1, // path length is 1 node
	}
	// r is a list of all nodes reached so far.
	// the chain of nodes following the prev member represents the
	// best path found so far from the start to this node.
	r := map[graph.EstimateNode]*rNode{start: p}
	// oh is a heap of nodes "open" for exploration.  nodes go on the heap
	// when they get an initial or new "g" path distance, and therefore a
	// new "f" which serves as priority for exploration.
	oh := openHeap{p}
	for len(oh) > 0 {
		bestPath := heap.Pop(&oh).(*rNode)
		bestNode := bestPath.nd
		if bestNode == end {
			// done
			dist := bestPath.g
			i := bestPath.n
			path := make([]graph.Neighbor, i)
			for i > 0 {
				i--
				path[i] = graph.Neighbor{bestPath.prevEdge, bestPath.nd}
				bestPath = bestPath.prevNode
			}
			return path, dist
		}
		bestNode.Visit(func(nb graph.Neighbor) {
			nd := nb.Nd.(graph.EstimateNode)
			ed := nb.Ed.(graph.DistanceEdge)
			g := bestPath.g + ed.Distance()
			if alt, reached := r[nd]; reached {
				if g > alt.g {
					// new path to nd is longer than some alternate path
					return
				}
				if g == alt.g && bestPath.n+1 >= alt.n {
					// new path has identical length of some alternate path
					// but it takes more hops.  go with fewest nodes in path.
					return
				}
				// cool, we found a better way to get to this node.
				// update alt with new data and make sure it's on the heap.
				alt.prevNode = bestPath
				alt.prevEdge = ed
				alt.g = g
				alt.f = g + alt.nd.Estimate(end)
				alt.n = bestPath.n + 1
				if alt.rx < 0 {
					heap.Push(&oh, alt)
				} else {
					heap.Fix(&oh, alt.rx)
				}
			} else {
				// bestNode being reached for the first time.
				p := &rNode{
					nd:       nd,
					prevNode: bestPath,
					prevEdge: ed,
					g:        g,
					f:        g + start.Estimate(end),
					n:        bestPath.n + 1,
				}
				r[nd] = p         // add to list of reached nodes
				heap.Push(&oh, p) // and it's now open for exploration
			}
		})
	}
	return nil, math.Inf(1) // no path
}

// AStarM is A* optimized for monotonic estimates.
//
// An admissable estimate may further be monotonic.  Monotonic means that if
// node B is a neighbor of node A with edge AB, then
// A.Estimate(C) <= AB.Distance() + B.Estimate(C).
func AStarM(start, end graph.EstimateNode) ([]graph.Neighbor, float64) {
	p := &rNode{
		nd: start,
		f:  start.Estimate(end),
		n:  1,
	}

	// difference from AStarA:
	// instead of r, a list of all nodes reached so far, there are two
	// lists, open and closed. open contains nodes "open" for exploration.
	// nodes are added to the list as they are reached, then moved to
	// closed as they are found to be on the best path.
	open := map[graph.EstimateNode]*rNode{start: p}
	closed := map[graph.EstimateNode]struct{}{}

	oh := openHeap{p}
	for len(oh) > 0 {
		bestPath := heap.Pop(&oh).(*rNode)
		bestNode := bestPath.nd
		if bestNode == end {
			// done
			dist := bestPath.g
			i := bestPath.n
			path := make([]graph.Neighbor, i)
			for bestPath != nil {
				i--
				path[i] = graph.Neighbor{bestPath.prevEdge, bestPath.nd}
				bestPath = bestPath.prevNode
			}
			return path, dist
		}

		// difference from AStarA:
		// move nodes to closed list as they are found to be best so far.
		delete(open, bestNode)
		closed[bestNode] = struct{}{}

		bestNode.Visit(func(nb graph.Neighbor) {
			nd := nb.Nd.(graph.EstimateNode)
			ed := nb.Ed.(graph.DistanceEdge)

			// difference from AStarA:
			// Monotonicity means that f cannot be improved.
			if _, ok := closed[nd]; ok {
				return
			}

			g := bestPath.g + ed.Distance()
			if alt, reached := open[nd]; reached {
				if g > alt.g {
					// new path to nd is longer than some alternate path
					return
				}
				if g == alt.g && bestPath.n+1 >= alt.n {
					// new path has identical length of some alternate path
					// but it takes more hops.  go with fewest nodes in path.
					return
				}
				// cool, we found a better way to get to this node.
				// update alt with new data and reheap.
				alt.prevNode = bestPath
				alt.prevEdge = ed
				alt.g = g
				alt.f = g + alt.nd.Estimate(end)
				alt.n = bestPath.n + 1

				// difference from AStarA:
				// we know alt was on the heap because we found it in the
				// open list.
				heap.Fix(&oh, alt.rx)
			} else {
				// bestNode being reached for the first time.
				p := &rNode{
					nd:       nd,
					prevNode: bestPath,
					prevEdge: ed,
					g:        g,
					f:        g + start.Estimate(end),
					n:        bestPath.n + 1,
				}
				open[nd] = p      // new node is now open for exploration.
				heap.Push(&oh, p) // keep heap matching open list.
			}
		})
	}
	return nil, math.Inf(1) // no path
}

// rNode holds data for a "reached" node
type rNode struct {
	nd       graph.EstimateNode
	prevNode *rNode             // chain encodes path back to start
	prevEdge graph.DistanceEdge // edge from prevNode to the node of this struct
	g        float64            // "g" best known path distance from start node
	f        float64            // "g+h", path dist + heuristic estimate
	n        int                // number of nodes in path
	rx       int                // heap.Remove index
}

type openHeap []*rNode

// implement container/heap
func (h openHeap) Len() int           { return len(h) }
func (h openHeap) Less(i, j int) bool { return h[i].f < h[j].f }
func (h openHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].rx = i
	h[j].rx = j
}
func (p *openHeap) Push(x interface{}) {
	h := *p
	rx := len(h)
	h = append(h, x.(*rNode))
	h[rx].rx = rx
	*p = h
}

func (p *openHeap) Pop() interface{} {
	h := *p
	last := len(h) - 1
	*p = h[:last]
	h[last].rx = -1
	return h[last]
}
