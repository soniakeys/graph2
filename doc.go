// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

// Graph implements graph algorithms Dijkstra's shortest path, algorithm A,
// and A*.
//
// Dijkstra's algorithm finds the shortest path between two nodes in a
// directed or undirected graph with weighted edges.  The edge weights
// must be non-negative.
//
// Algorithm A and A* optimize the shortest path search using a heuristic
// estimate of the distance to the end node.  A heuristic is termed admissable
// if the estimate is always less than or equal to the actual path distance.
// In this case he algorithm is guaranteed to find the shortest path and is
// termed A*.  If the heuristic is inadmissable, the same algorithm will still
// find a path but possibly not the shortest path.  In this case the algorithm
// is termed algorithm A.
//
// An admissable estimate may further be monotonic.  Monotonic means that if
// node B is a neighbor of node A with edge AB, then the estimate from A
// must be less than or equal to the edge distance AB plus the estimate from
// B.  The package has a separate function optimized for monotonic graphs.
//
// Graph requires Go 1.2.  (It uses the new heap.Fix method.)
package graph
