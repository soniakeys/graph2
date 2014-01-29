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

// Neighbor associates an arc or edge with the node that is reached by the
// arc or edge.
type Neighbor struct {
	Ed interface{} // arc or edge
	Nd NeighborNode
}

// Arc and Edge are completely generic to hold any object but are separate
// named types to indicate interpretation.
type (
	Arc  interface{} // directed
	Edge interface{} // undirected
)

// Weighted is an object such as an arc or edge that describes a weight,
// typically a non-negative quantity.  Alternate terms for weight include
// distance, length, and cost.
type Weighted interface {
	Weight() float64
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

// ArborNode enables construction of an arborescence, or directed spanning tree.
//
// A function that traverses a graph can use this interface to construct
// an arborescence on top of the graph.  The function should call LinkFrom
// for each node of the graph.  The call for the first node should have
// nil arguments.  The caller should retain the result from this first
// call as the root of the arborescence.  Subsequent calls pass a value
// previously returned by LinkFrom as prev, and the connecting edge from
// the original graph as arc.
//
// Note that implementations of ArborNode also determine the implementation
// of the NeighborNode returned by LinkFrom.  The two types need not be
// the same.
type ArborNode interface {
	NeighborNode
	// LinkFrom should construct a new node based on the reciever and
	// construct a link based on arc that links prev to the new node.
	LinkFrom(prev NeighborNode, arc Arc) NeighborNode
}

// SpannerNode enables construction of a spanning tree.
//
// A function that traverses a graph can use this interface to construct
// a spanning tree on top of the graph.  The function should call Span
// for each node of the graph.  The call for the first node should have
// nil arguments.  Subsequent calls pass a value previously returned by
// Span as prev, and the connecting edge from the original graph as ed.
//
// Note that implementations of SpannerNode also determine the implementation
// of the NeighborNode returned by Span.  The two types need not be the same.
type SpannerNode interface {
	NeighborNode
	// Span should construct a new node based on the reciever and
	// construct a link based on ed that links prev to the new node.
	Span(prev NeighborNode, ed Edge) NeighborNode
}
