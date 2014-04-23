// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

// Graph2 defines interfaces and other types useful for graph algorithms.
//
// Subdirectory search contains graph search functions.  Implemented search
// algorithms are Dijkstra’s shortest path, A*, and algorithm A.  Functions
// in package search operate through the interfaces in package graph and make
// no requirements about concrete graph representation.
//
// Subdirectory adj contains concrete types and methods for an adjacency graph.
// The types are adequate for exercising the functions in package search and
// are generalized to be useful for other applications.
//
// Neither search nor adj depend on the other; they only depend on graph.
package graph2
