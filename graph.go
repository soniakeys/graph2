// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph

// graph.go has interfaces that need to be implemented on a graph
// representation for various graph search functions in this package.

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
	DistanceNeighbors([]DistanceNeighbor) []DistanceNeighbor
}

// DistanceEdge is an edge that describes a distance, typically a
// non-negative quantity.  Some graph search algorithms require non-negative
// edge distances.  Check the documentation for the graph search function
// you will use.
//
// The only data discribed is the distance along the edge.
// Alternate terms for distance include length, weight, and cost.
type DistanceEdge interface {
	// Distance returns distance along edge as a float64.
	// The value should not be an Inf or NaN.
	Distance() float64
}

// DistanceNeighbor describes a distance relationship of one node to another.
// DistanceNeightbor specifies two things, an edge leading from a node, and
// the node at the other end of the edge.  It thus represents a directed link.
// In an undirected graph constructed with DistanceNeighbors, a neighbor of
// one node always has the first node as one of its neighbors and further,
// the distance is the same in each direction.
type DistanceNeighbor struct {
	DistanceEdge
	DistanceNode
}

// EstimateNode describes a node that can provide a distance estimate
// to another node.  This estimate is often called h, or a heuristic estimate,
// as it has heuristic use.  Admissability and monoticity are unspecified.
//
// Admissable means the value returned by Estimate must be less than or equal
// to the actual path distance.
//
// An admissable estimate may further be monotonic.  Monotonic means that if
// node B is a neighbor of node A with edge AB, then
// A.Estimate(C) <= AB.Distance() + B.Estimate(C).
//
// Some graph search algorithms require admissability or monotonicity.
// Check the documentation for the graph search function you will use.
type EstimateNode interface {
	// EstimateNeighbors works like DistanceNode.DistanceNeighbors.
	// See DistanceNode documentation.
	EstimateNeighbors([]EstimateNeighbor) []EstimateNeighbor
	Estimate(EstimateNode) float64
}

// EstimateNeighbor describes a distance relationship of one node to another.
// EstimateNeightbor specifies two things, an edge leading from a node, and
// the node at the other end of the edge.  It thus represents a directed link.
// In an undirected graph constructed with EstimateNeighbors, a neighbor of
// one node always has the first node as one of its neighbors and further,
// the distance is the same in each direction.
type EstimateNeighbor struct {
	DistanceEdge
	EstimateNode
}
