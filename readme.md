Graph
=====

Graph implements Dijkstra’s shortest path algorithm.

[![Build Status](https://travis-ci.org/soniakeys/graph.png)](https://travis-ci.org/soniakeys/graph)

[![GoDoc](https://godoc.org/github.com/garyburd/gddo?status.png)](http://godoc.org/github.com/soniakeys/graph)

[![Go Walker](http://gowalker.org/api/v1/badge)](http://gowalker.org/github.com/soniakeys/graph)

The package provides a way to do shortest path searches on existing data
structures.  Rather than specify types you must use for graphs, nodes, and
edges, the package specifies some interfaces you must implement on your
existing types.  You are not required to host any of the bookkeeping data
needed by the search algorithm.  The methods you implement simply return
information about your graph.

Multiple searches can run concurrently on the same graph.
