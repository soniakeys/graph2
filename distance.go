// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph

// DistanceNode describes a node (or vertex) with a distance relationship to
// other nodes.  Nodes to which it has a direct distance relationship are
// termed its "neighbors."
type DistanceNode interface {
	// Neighbors returns neighbors of this node; that is, nodes reachable
	// from this one by following a single edge.  For an undirected graph
	// constructed with DistanceNodes, if a node B is in the list returned
	// by A.Neighbors() then A must be in the list returned by B.Neighbors().
	// Further, the edge returned in each case must represent the same
	// distance, by the result of edge.Distance().
	//
	// The argument to Neighbors must be a zero length slice.  It can be
	// nil or a slice with zero length and non-zero capacity.  This is
	// useful in cases where the implementation of DistanceNode constructs
	// the neighbor list on demand.  Functions that call Neighbors
	// (e.g. DijkstraShortestPath) can pass the the previous result truncated
	// to zero length.  Implementations can append to this slice, reusing the
	// existing capacity and minimizing garbage.
	//
	// For performance concerns, an implentation might maintain a neighbor
	// list for each node as part of its graph representation and simply
	// return the list.  In this case implementations will ignore the
	// argument to Neighbors.
	Neighbors([]DistanceNeighbor) []DistanceNeighbor
}

// DistanceEdge is an edge that describes a distance, or non-negative quantity.
// The only data discribed is the distance along the edge.
// Alternate terms for distance include length, weight, and cost.
type DistanceEdge interface {
	// Distance returns distance along edge as a float64.
	// The value should be non-negative, not an Inf, and not a NaN.
	Distance() float64
}

// DistanceNeighbor describes a distance relationship of one node to another,
// distance being non-negative.  DistanceNeightbor specifies two things,
// an edge leading from a node, and the node at the other end of the edge.
// It thus represents a directed link.  In an undirected graph constructed
// with DistanceNeighbors, a neighbor of one node always has the first node
// as one of its neighbors and further, the distance is the same in each
// direction.
type DistanceNeighbor struct {
	DistanceEdge
	DistanceNode
}
