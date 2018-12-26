package state

import (
	"fmt"
	"log"
)

type Node struct {
	Planet     *Planet
	id         int
	former     int
	dist       int
	generation uint32
}

func (a Node) Equal(b *Node) bool {
	return a.Planet == b.Planet
}

type Graph struct {
	nodes            []Node
	translationTable map[*Planet]*Node
	generation       uint32
}

type edge struct {
	from int
	to   int
	dist float64
}

type Edges []edge

func (a Edges) Len() int      { return len(a) }
func (a Edges) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a Edges) Less(i, j int) bool {
	return a[i].dist < a[j].dist
}

func NewGraph(planets []*Planet) Graph {

	g := Graph{
		nodes:            make([]Node, len(planets)),
		translationTable: make(map[*Planet]*Node),
		generation:       0,
	}
	for i, planet := range planets {
		g.nodes[i].id = i
		g.nodes[i].Planet = planet
		g.translationTable[planet] = &g.nodes[i]
	}

	return g
}

func (g *Graph) Visit(planet *Planet, f func(n Node) bool) {
	g.generation += 2
	queue := make([]*Node, 0)

	queue = append(queue, g.translationTable[planet])
	for len(queue) > 0 {
		elem := queue[0]
		queue = queue[1:]

		cont := f(g.nodes[elem.id])
		elem.generation = g.generation
		if !cont {
			return
		}

		// visit children
		for _, child := range elem.Planet.Connected {
			if g.translationTable[child].generation == g.generation {
				// already visited
				continue
			}
			queue = append(queue, g.translationTable[child])

			if g.translationTable[child].generation < g.generation-1 || // never seen this child
				g.translationTable[child].dist > elem.dist+1 { // the current path is shorter than the old one

				// set distance and path
				g.translationTable[child].dist = elem.dist + 1
				g.translationTable[child].former = elem.id
				// mark as seen
				g.translationTable[child].generation = g.generation - 1
			}
		}
	}
}

func (g *Graph) ShortestPath(root *Planet, target *Planet, graphIsPrepared bool) []*Planet {
	if ! graphIsPrepared {
		g.Visit(root, func(n Node) bool {
			return !g.translationTable[target].Equal(&n)
		})
	}

	path := make([]*Planet, 0)

	cur := *g.translationTable[target]
	path = append(path, target)
	for !cur.Equal(g.translationTable[root]) {
		path = append(path, g.nodes[cur.former].Planet)
		cur = g.nodes[cur.former]
		if len(path) > 500 {
			fmt.Println(root, g.translationTable[root].Planet.Connected)
			fmt.Println(target, g.translationTable[target].Planet.Connected)
			fmt.Println(path)
			for i := 0; i < 5; i++ {
				fmt.Println("  ", path[i], g.translationTable[path[i]].Planet.Connected, "dist", g.translationTable[path[i]].dist, "former", g.translationTable[path[i]].former, "gen", g.translationTable[path[i]].generation)
			}
			log.Panicf("Unintended behavior when checking path from %v to %v, %v", root, target, path)
		}
	}

	//revert path
	for i := 0; i < len(path)/2; i++ {
		path[i], path[len(path)-1-i ] = path[len(path)-1-i ], path[i]
	}

	return path
}

func (g *Graph) GraphSize(root Node) int {

	count := 0

	g.Visit(root.Planet, func(planet Node) bool {
		count++
		return true
	})

	return count
}

func (g *Graph) HasCycle(root *Planet) bool {
	g.generation++

	// remember which edges we already traveled
	usedEdges := make(map[edge]bool)

	// queue for bfs
	queue := make([]*Node, 0)
	queue = append(queue, g.translationTable[root])

	for len(queue) > 0 {

		elem := queue[0]
		queue = queue[1:]

		// already visited
		if elem.generation == g.generation {
			return true
		}

		elem.generation = g.generation

		// visit children
		for _, child := range elem.Planet.Connected {
			e := edge{elem.id, g.translationTable[child].id, 0}
			// we come by this edge to this node therefore we can ignore it
			if usedEdges[e] {
				continue
			}
			usedEdges[edge{g.translationTable[child].id, elem.id, 0}] = true
			queue = append(queue, g.translationTable[child])
		}
	}

	return false
}
