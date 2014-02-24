package search

import "github.com/soniakeys/graph"

// DepthFirst traverses nodes in depth first order.
//
// The visitor function is called for each node, starting with the argument n.
// If the visitor function returns false for any node, DepthFirst stops and
// returns false immediately.  DepthFirst returns true otherwise.
func DepthFirst(n graph.Node, v graph.LevelVisitor) (ok bool) {
	m := map[graph.Node]struct{}{}
	var r func(graph.Node, int) bool
	r = func(n graph.Node, level int) bool {
		if _, ok := m[n]; ok {
			return true
		}
		if !v(n, level) {
			return false
		}
		m[n] = struct{}{}
		level++
		return n.VisitAdjNodes(func(n graph.Node) bool {
			return r(n, level)
		})
	}
	return r(n, 0)
}
