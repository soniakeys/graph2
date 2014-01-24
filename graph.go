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

// SpannerNode enables construction of a spanning tree.
//
// A function that traverses a graph can use this interface to construct
// a spanning tree on top of the graph.  The function should call LinkFrom
// for each node of the graph.  The call for the first node should have
// nil arguments.  Subsequent calls pass a value previously returned by
// LinkFrom as prev, and the connecting edge from the original graph as ed.
//
// Note that implentations of SpannerNode also determine the implentation
// of the NeighborNode returned by LinkFrom.  The two types need not be
// the same.
type SpannerNode interface {
	NeighborNode
	// LinkFrom should construct a new node based on the reciever and
	// construct a link based on ed that links prev to the new node.
	LinkFrom(prev NeighborNode, ed Edge) NeighborNode
}
