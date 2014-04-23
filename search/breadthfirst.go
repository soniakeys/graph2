package search

import "github.com/soniakeys/graph2"

func BreadthFirst2(g graph2.BF2Graph, start graph2.BF2Node, visit graph2.BF2NodeVisitor) (p map[graph2.BF2Node]graph2.BF2Node, ok bool) {
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
		level = make([]graph2.BF2Node, len(lMap))
		i := 0
		for n := range lMap {
			level[i] = n
			i++
		}
	}
	return visited, true
}

type bfMap map[graph2.BF2Node]struct{}
type bfVMap map[graph2.BF2Node]graph2.BF2Node
type bfList []graph2.BF2Node

func topDown(lNum int, level []graph2.BF2Node, visit graph2.BF2NodeVisitor, visited bfVMap, unvis bfMap) (next bfList, mnext int, ok bool) {
	for _, v := range level {
		if !v.VisitBF2Out(func(n graph2.BF2Node) int {
			if _, ok := visited[n]; ok {
				return graph2.BF2Go
			}
			if !visit(n, lNum) {
				return graph2.BF2Stop
			}
			visited[n] = v
			delete(unvis, n)
			next = append(next, n)
			mnext += n.NumAdj()
			return graph2.BF2Go
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

func bottomUp(lNum int, lmap bfMap, visit graph2.BF2NodeVisitor, visited bfVMap, unvis bfMap) (next bfMap, mnext int, ok bool) {
	next = bfMap{}
	for v := range unvis {
		if !v.VisitBF2In(func(n graph2.BF2Node) int {
			if _, ok := lmap[n]; !ok {
				return graph2.BF2Go
			}
			if !visit(v, lNum) {
				return graph2.BF2Stop
			}
			visited[v] = n
			delete(unvis, v)
			next[v] = struct{}{}
			mnext += v.NumAdj()
			return graph2.BF2Found
		}) {
			return nil, 0, false
		}
	}
	return next, mnext, true
}

func BreadthFirst1(start graph2.Node, visit graph2.LevelVisitor) (p map[graph2.Node]graph2.Node, ok bool) {
	lnum := 0
	visited := map[graph2.Node]graph2.Node{}
	if !visit(start, lnum) {
		return visited, false
	}
	visited[start] = nil
	level := []graph2.Node{start}
	next := []graph2.Node{}
	for len(level) > 0 {
		lnum++
		for _, v := range level {
			if !v.VisitAdjNodes(func(n graph2.Node) bool {
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
