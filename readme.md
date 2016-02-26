# Graph2

Graph2 defines interfaces and other types useful for graph algorithms.

This was an experiment in minimizing interfaces.  (It is no longer under
development.)  If the interface for a node
presents the node’s relationship to other nodes, that is enough for some
graph algorithms.  That is, no named type is needed for the graph as a whole.
In general, a function implementing a graph algorithm should have parameters
with method sets no larger than what the algorithm really needs.

Work on this package slowed at some point as I began to put more effort into
another approach--not having interfaces at all!
See github.com/soniakeys/graph for this approach.

Further to this package though, subdirectory search contains graph search
and traversal functions.  Implemented algorithms are Dijkstra’s shortest path,
A\*, algorithm A, depth first, breadth first, and Beamer’s direction-optimizing
breadth first.

Subdirectory adj contains concrete types and methods for an adjacency list
graph representation.
