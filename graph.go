// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph

// A Node represents an adjacency relationship to other nodes.
//
// The relationship defines arcs or edges but does not associate any
// data objects with the arcs or edges.  A graph of Nodes is necessarily
// unweighted then.
type Node interface {
	// Visit calls the NodeVisitor function for neighboring nodes.
	// If the NodeVisitor function returns false for any neighbor, Visit
	// stops visiting and returns false immediately.  Visit returns true
	// otherwise.
	VisitOk(NodeOkVisitor) (ok bool)
}

// NodeVisitor is the argument type for Node.Visit.  A node visitor
// can process the node as appropriate and return a result to the caller
// indicating whether to continue visiting or not.  True means continue,
// false means there is no need to continue.
type NodeOkVisitor func(Node) (ok bool)

// BFGraph must be implemented on a collection of nodes for the BreadthFirst
// search algorithm.
type BFGraph interface {
	// Nodes must construct a new map populated with all nodes in the graph.
	// BreadthFirst2 will consume the map.
	Nodes() map[BFNode]struct{}
	// Total number of edges or arcs in the graph.  For an undirected graph
	// return the number of graph edges.  While each edge may be represented
	// internally with links in both directions, count each edge once and do
	// not count the links separately.  For directional graphs, return the
	// number of graph arcs.  Again, if arcs are duplicated internally for
	// access as both inward and outward pointing arcs, count each arc only
	// once.
	NumEdges() int
}

// BFNeighborVisitors are defined within the BreadthFirst search function.
// Neither graph implementers nor search clients need to define these.
// BreadthFirst passes them as arguments to the Visit methods of BFNode.
type BFNeighborVisitor func(BFNode) int

// BF constants returned by BFNeighborVisitor functions.  They communicate
// a result from the BreadthFirst search function to (BFNode) methods
// implementing a breadth first searchable graph.
const (
	BFGo    = iota // Continue visiting neighbors.
	BFStop         // Stop visiting, return false to signal stop bf search.
	BFFound        // Stop visiting, return true.
)

// BFNode defines methods that the BreadthFirst search function will call
// in the course of its search.  The VisitBF functions iterate over node
// neighbors.  A false result from any of these tells BreadthFirst to terminate
// search early.
type BFNode interface {
	// VisitBFIn must iterate over the inward-pointing arcs of the node.
	// for each neighbor by inward-pointing arc, VisithBFIn must call the
	// neighbor visitor function and handle one of three integer results
	// as follows:
	//    BFGo:    Continue iteration.
	//    BFStop:  Break from iteration and return false from VisitBFIn.
	//    BFFound: Break from iteration but still return true from VisitBFIn.
	// If the node has no neighbors by inward-pointing arcs or if the visitor
	// function returns BFGo for all neighbors, VisitBFIn must return true.
	VisitBFIn(BFNeighborVisitor) (ok bool)
	// VisitBFOut must iterate over the outward-pointing arcs of the node,
	// The integer result of the visitor function must be handled as for
	// VisitBFIn except that the only possible results are BFGo and BFStop.
	VisitBFOut(BFNeighborVisitor) (ok bool)
	// NumAdj must return the number of outward-pointing arcs from the node.
	NumAdj() int
}

// A BFNodeVisitor is implemented by a caller of the BreadthFirst search
// function.  BreadthFirst will call the BFNodeVisitor for each node of a
// graph as it traverses the graph in breadth first order.  Argument n is
// the node being visited, level is the level of the search where 0 is the
// start node, 1 is immediate neighbors of the node, and so on.  BFNodeVistor
// should return true for BreadthFirst to continue traversing the graph.
// It can return false to signal BreadthFirst to terminate traversal early.
type BFNodeVisitor func(n BFNode, level int) (ok bool)

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
