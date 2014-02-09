package search

import "github.com/soniakeys/graph"

func BreadthFirst(start graph.Node, visit graph.NodeOkVisitor) (ok bool) {
	if !visit(start) {
		return false
	}
	visited := map[graph.Node]struct{}{start: struct{}{}}
	level := []graph.Node{start}
	next := []graph.Node{}
	for len(level) > 0 {
		for _, v := range level {
			if !v.VisitOk(func(n graph.Node) bool {
				if _, ok := visited[n]; !ok {
					if !visit(n) {
						return false
					}
					visited[n] = struct{}{}
					next = append(next, n)
				}
				return true
			}) {
				return false
			}
		}
		level = next
		next = level[:0]
	}
	return true
}
