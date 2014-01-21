// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph

// A NeighborNode has some way of visiting its neighbors.
//
// Visit should call the NeighborVisitor function for each neighbor
// of the receiver.
type NeighborNode interface {
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
// non-negative quantity.  Alternate terms for distance include length,
// weight, and cost.
type DistanceEdge interface {
	Distance() float64
}

// An Estimator provides a distance estimate from itself to an EstimateNode.
// This estimate is often called h, or a heuristic estimate.
type Estimator interface {
	Estimate(EstimateNode) float64
}

// EstimateNode describes a node that can provide a distance estimate
// to another EstimateNode.
type EstimateNode interface {
	NeighborNode
	Estimator
}

// Spanner enables construction of a spanning tree.
type Spanner interface {
	NeighborNode
	// LinkFrom should construct a new node based on the reciever and
	// construct a link based on ed that links prev to the new node.
	LinkFrom(prev NeighborNode, ed Edge) NeighborNode
}
