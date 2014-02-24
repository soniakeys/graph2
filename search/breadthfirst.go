package search

import "github.com/soniakeys/graph"

func BreadthFirst2(g graph.BF2Graph, start graph.BF2Node, visit graph.BF2NodeVisitor) (p map[graph.BF2Node]graph.BF2Node, ok bool) {
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
		// TODO mf not needed here.  compute on conversion back to top down.
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
		level = make([]graph.BF2Node, len(lMap))
		i := 0
		for n := range lMap {
			level[i] = n
			i++
		}
	}
	return visited, true
}

type bfMap map[graph.BF2Node]struct{}
type bfVMap map[graph.BF2Node]graph.BF2Node
type bfList []graph.BF2Node

func topDown(lNum int, level []graph.BF2Node, visit graph.BF2NodeVisitor, visited bfVMap, unvis bfMap) (next bfList, mnext int, ok bool) {
	for _, v := range level {
		if !v.VisitBF2Out(func(n graph.BF2Node) int {
			if _, ok := visited[n]; ok {
				return graph.BF2Go
			}
			if !visit(n, lNum) {
				return graph.BF2Stop
			}
			visited[n] = v
			delete(unvis, n)
			next = append(next, n)
			mnext += n.NumAdj()
			return graph.BF2Go
		}) {
			return nil, 0, false
		}
	}
	return next, mnext, true
}

const (
	bf2Go    = iota // continue searching neighbors
	bf2Stop         // stop searching neighbors, stop bf search
	bf2Found        // stop searching neigbhors, continue bf search
)

func bottomUp(lNum int, lmap bfMap, visit graph.BF2NodeVisitor, visited bfVMap, unvis bfMap) (next bfMap, mnext int, ok bool) {
	next = bfMap{}
	for v := range unvis {
		if !v.VisitBF2In(func(n graph.BF2Node) int {
			if _, ok := lmap[n]; !ok {
				return graph.BF2Go
			}
			if !visit(v, lNum) {
				return graph.BF2Stop
			}
			visited[v] = n
			delete(unvis, v)
			next[v] = struct{}{}
			mnext += v.NumAdj()
			return graph.BF2Found
		}) {
			return nil, 0, false
		}
	}
	return next, mnext, true
}

func BreadthFirst1(start graph.Node, visit graph.LevelVisitor) (p map[graph.Node]graph.Node, ok bool) {
	lnum := 0
	visited := map[graph.Node]graph.Node{}
	if !visit(start, lnum) {
		return visited, false
	}
	visited[start] = nil
	level := []graph.Node{start}
	next := []graph.Node{}
	for len(level) > 0 {
		lnum++
		for _, v := range level {
			if !v.VisitAdjNodes(func(n graph.Node) bool {
				if _, ok := visited[n]; !ok {
					if !visit(n, lnum) {
						return false
					}
					visited[n] = v
					next = append(next, n)
				}
				return true
			}) {
				return visited, false
			}
		}
		level = next
		next = level[:0]
	}
	return visited, true
}
