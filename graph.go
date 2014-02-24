// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph

// A Node represents an adjacency relationship to other nodes
//
// The relationship defines arcs or edges but does not associate any
// data objects with the arcs or edges.  A graph of Nodes is necessarily
// unweighted then.
type Node interface {
	// VisitAdjNodes must iterate over nodes adjacent by outward pointing arcs
	// or undirected edges and call the visitor function for each.  If the
	// visitor function returns false for any adjacent node, VisitAdjNode
	// should stop visiting and return false immediately.  VisitAdjNode should
	// return true otherwise.
	VisitAdjNodes(AdjNodeVisitor) (ok bool)
}

// An AdjNodeVisitor is defined within a search or traversal function.  Neither
// graph implementers nor search clients need to define one.  The search or
// traversal function passes it as the argument to the VisitAdjNodes method of
// Node.
type AdjNodeVisitor func(n Node) (ok bool)

// A LevelVisitor is an argument to some search or traversal functions.
// The search or traversal function will call LevelVisitor for each node of a
// graph as the function traverses the graph.  Argument n is the node being
// visited, level is the level of the search--the number of edges in the
// path from the start node to n.  Thus level will be 0 for the start node,
// 1 for nodes immediately adjacent, and so on. LevelVistor should return true
// to indicate that the search function should continue traversing the graph.
// LevelVisitor can return false to indicate that search can be terminated
// immediately.
type LevelVisitor func(n Node, level int) (ok bool)

// BF2Graph must be implemented on a collection of nodes for the BreadthFirst2
// search algorithm.
type BF2Graph interface {
	// Nodes must construct a new map populated with all nodes in the graph.
	// BreadthFirst2 will consume the map.
	Nodes() map[BF2Node]struct{}
	// Total number of edges or arcs in the graph.  For an undirected graph
	// return the number of graph edges.  While each edge may be represented
	// internally with links in both directions, count each edge once and do
	// not count the links separately.  For directional graphs, return the
	// number of graph arcs.  Again, if arcs are duplicated internally for
	// access as both inward and outward pointing arcs, count each arc only
	// once.
	NumEdges() int
}

// BF2NeighborVisitors are defined within the BreadthFirst2 search function.
// Neither graph implementers nor search clients need to define these.
// BreadthFirst2 passes them as arguments to the Visit methods of BF2Node.
type BF2NeighborVisitor func(BF2Node) int

// BF2 constants returned by BF2NeighborVisitor functions.  They communicate
// a result from the BreadthFirst2 search function to (BF2Node) methods
// implementing a BF2 searchable graph.
const (
	BF2Go    = iota // Continue visiting neighbors.
	BF2Stop         // Stop visiting, return false to signal stop bf search.
	BF2Found        // Stop visiting, return true.
)

// BF2Node defines methods that the BreadthFirst2 search function will call
// in the course of its search.  The VisitBF2 functions iterate over node
// neighbors.  A false result from any of these tells BreadthFirst2 to terminate
// search early.
type BF2Node interface {
	// VisitBF2In must iterate over the inward-pointing arcs of the node.
	// for each neighbor by inward-pointing arc, VisithBF2In must call the
	// neighbor visitor function and handle one of three integer results
	// as follows:
	//    BF2Go:    Continue iteration.
	//    BF2Stop:  Break from iteration and return false from VisitBF2In.
	//    BF2Found: Break from iteration but still return true from VisitBF2In.
	// If the node has no neighbors by inward-pointing arcs or if the visitor
	// function returns BF2Go for all neighbors, VisitBF2In must return true.
	VisitBF2In(BF2NeighborVisitor) (ok bool)
	// VisitBF2Out must iterate over the outward-pointing arcs of the node,
	// The integer result of the visitor function must be handled as for
	// VisitBF2In except that the only possible results are BF2Go and BF2Stop.
	// BF2Found does not need to be handled.
	VisitBF2Out(BF2NeighborVisitor) (ok bool)
	// NumAdj must return the number of outward-pointing arcs from the node.
	NumAdj() int
}

// A BF2NodeVisitor is implemented by a caller of the BreadthFirst2 search
// function.  BreadthFirst2 will call the BF2NodeVisitor for each node of a
// graph as it traverses the graph in breadth first order.  Argument n is
// the node being visited, level is the level of the search where 0 is the
// start node, 1 is immediate neighbors of the node, and so on.  BF2NodeVistor
// should return true for BreadthFirst2 to continue traversing the graph.
// It can return false to signal BreadthFirst2 to terminate traversal early.
type BF2NodeVisitor func(n BF2Node, level int) (ok bool)

// A HalfNode represents an adjacency relationship.
//
// The relationship is by half arcs or half edges that directly connect to
// other HalfNodes.
type HalfNode interface {
	// Visit should call the AdjHalfVisitor function for each adjacent half
	// arc or half edge.
	VisitAdjHalfs(AdjHalfVisitor)
}

// AdjHalfVisitor is the argument type for HalfNode.VisitAdjHalfs.
type AdjHalfVisitor func(Half)

// Half is a half arc or half edge.  It associates an arc or edge with the
// node that is reached by the arc or edge.  For a node in a directed graph,
// a Half is an arc leading from the node and the node at the end of the arc.
// For a node in an undirected graph, a Half is an edge touching the node and
// the node on the other end of the edge.
type Half struct {
	Ed interface{} // arc or edge
	Nd HalfNode
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
	HalfNode
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
// of the HalfNode returned by LinkFrom.  The two types need not be
// the same.
type ArborNode interface {
	HalfNode
	// LinkFrom should construct a new node based on the reciever and
	// construct a link based on arc that links prev to the new node.
	LinkFrom(prev HalfNode, arc Arc) HalfNode
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
// of the HalfNode returned by Span.  The two types need not be the same.
type SpannerNode interface {
	HalfNode
	// Span should construct a new node based on the reciever and
	// construct a link based on ed that links prev to the new node.
	Span(prev HalfNode, ed Edge) HalfNode
}
