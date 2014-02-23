package search

import "github.com/soniakeys/graph"

func BreadthFirst(g graph.BFGraph, start graph.BFNode, visit graph.BFNodeVisitor) (p map[graph.BFNode]graph.BFNode, ok bool) {
	lNum := 0
	visited := bfVMap{}
	if !visit(start, lNum) {
		return visited, false
	}
	visited[start] = nil
	level := bfList{start}
	mf := start.NumAdj()
	m := g.NumEdges()
	ctb := m / 10
	unvis := g.Nodes()
	n := len(unvis)
	k14 := 14 * m / n // 14 * mean degree
	cbt := n / k14
	delete(unvis, start)
	for {
		lNum++
		level, mf, ok = topDown(lNum, level, visit, visited, unvis)
		if !ok {
			return visited, false
		}
		if len(level) == 0 {
			return visited, true
		}
		if mf > ctb {
			// switch to bottom up!
		} else {
			// stick with top down
			continue
		}
		// convert
		lMap := make(bfMap, len(level))
		for _, n := range level {
			lMap[n] = struct{}{}
		}
	bottomUpAgain:
		lNum++
		lMap, mf, ok = bottomUp(lNum, lMap, visit, visited, unvis)
		if !ok {
			return visited, false
		}
		if len(lMap) == 0 {
			return visited, true
		}
		if len(lMap) < cbt {
			// switch back to top down!
		} else {
			// stick with bottom up
			goto bottomUpAgain
		}
		// convert
		level = make([]graph.BFNode, len(lMap))
		i := 0
		for n := range lMap {
			level[i] = n
			i++
		}
	}
	return visited, true
}

type bfMap map[graph.BFNode]struct{}
type bfVMap map[graph.BFNode]graph.BFNode
type bfList []graph.BFNode

func topDown(lNum int, level []graph.BFNode, visit graph.BFNodeVisitor, visited bfVMap, unvis bfMap) (next bfList, mnext int, ok bool) {
	for _, v := range level {
		if !v.VisitBFOut(func(n graph.BFNode) int {
			if _, ok := visited[n]; ok {
				return graph.BFGo
			}
			if !visit(n, lNum) {
				return graph.BFStop
			}
			visited[n] = v
			delete(unvis, n)
			next = append(next, n)
			mnext += n.NumAdj()
			return graph.BFGo
		}) {
			return nil, 0, false
		}
	}
	return next, mnext, true
}

const (
	bfGo    = iota // continue searching neighbors
	bfStop         // stop searching neighbors, stop bf search
	bfFound        // stop searching neigbhors, continue bf search
)

func bottomUp(lNum int, lmap bfMap, visit graph.BFNodeVisitor, visited bfVMap, unvis bfMap) (next bfMap, mnext int, ok bool) {
	next = bfMap{}
	for v := range unvis {
		if !v.VisitBFIn(func(n graph.BFNode) int {
			if _, ok := lmap[n]; !ok {
				return graph.BFGo
			}
			if !visit(v, lNum) {
				return graph.BFStop
			}
			visited[v] = n
			delete(unvis, v)
			next[v] = struct{}{}
			mnext += v.NumAdj()
			return graph.BFFound
		}) {
			return nil, 0, false
		}
	}
	return next, mnext, true
}

func BreadthFirstSimple(start graph.Node, visit graph.NodeOkVisitor) (ok bool) {
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
