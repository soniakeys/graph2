package search

import "github.com/soniakeys/graph2"

// DepthFirst traverses nodes in depth first order.
//
// The visitor function is called for each node, starting with the argument n.
// If the visitor function returns false for any node, DepthFirst stops and
// returns false immediately.  DepthFirst returns true otherwise.
func DepthFirst(n graph2.Node, v graph2.LevelVisitor) (ok bool) {
	m := map[graph2.Node]struct{}{}
	var r func(graph2.Node, int) bool
	r = func(n graph2.Node, level int) bool {
		if _, ok := m[n]; ok {
			return true
		}
		if !v(n, level) {
			return false
		}
		m[n] = struct{}{}
		level++
		return n.VisitAdjNodes(func(n graph2.Node) bool {
			return r(n, level)
		})
	}
	return r(n, 0)
}
