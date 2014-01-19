// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph

// graph.go has interfaces that need to be implemented on a graph
// representation for various graph search functions in this package.

// A NeighborNode has some way of visiting it's neighbors.
type NeighborNode interface {
	// Visit should call the NeighborVisitor function for each neighbor
	// of the receiver.
	Visit(NeighborVisitor)
}

// An algorithm can process neighbors of a NeighborNode by passing a
// NeighborVisitor to NeigborNode.Visit.
type NeighborVisitor func(Neighbor)

// Neighbor associates an edge with the node that is reached by the edge.
type Neighbor struct {
	Ed Edge
	Nd NeighborNode
}

// Edge is completely generic to hold any object representing an edge.
type Edge interface{}

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
	NeighborNode
	Estimator
}
type Estimator interface {
	Estimate(EstimateNode) float64
}
