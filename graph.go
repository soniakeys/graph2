// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph

type Node interface {
	Visit(NodeVisitor) (ok bool)
}

type NodeVisitor func(Node) (ok bool)

// An AdjNode represents an adjacency relationship.
//
// The relationship is by edges or arcs that directly connect to other
// AdjNodes.
type AdjNode interface {
	// Visit should call the HalfVisitor function for each adjacent half
	// arc or half edge.
	VisitAdj(HalfVisitor)
}

// HalfVisitor is the argument type for AdjNode.Visit.
type HalfVisitor func(Half)

// Half is a half edge or half arc.  It associates an arc or edge with the
// node that is reached by the arc or edge.  For a node in a directed graph,
// a Half is an arc leading from the node and the node at the end of the arc.
// For a node in an undirected graph, a Half is an edge touching the node and
// the node on the other end of the edge.
type Half struct {
	Ed interface{} // arc or edge
	Nd AdjNode
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
	AdjNode
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
// of the AdjNode returned by LinkFrom.  The two types need not be
// the same.
type ArborNode interface {
	AdjNode
	// LinkFrom should construct a new node based on the reciever and
	// construct a link based on arc that links prev to the new node.
	LinkFrom(prev AdjNode, arc Arc) AdjNode
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
// of the AdjNode returned by Span.  The two types need not be the same.
type SpannerNode interface {
	AdjNode
	// Span should construct a new node based on the reciever and
	// construct a link based on ed that links prev to the new node.
	Span(prev AdjNode, ed Edge) AdjNode
}
