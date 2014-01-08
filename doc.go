// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

// Graph implements graph algorithms (currently just Dijkstra's shortest path
// algorithm.)
//
// Dijkstra's algorithm is efficient for finding the shortest path between
// two nodes in a general directed or undirected graph.  The path length
// minimized is the sum of edge lengths in the path, which must be non-negative.
//
// Graph requires Go 1.2.  (It uses the new heap.Fix method.)
package graph
