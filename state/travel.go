package state

import (
	pb "github.com/eqinox76/RiseAndFallOfEmpires/proto"
	"fmt"
	"log"
)

type Node struct {
	Planet     *pb.Planet
	former     uint32
	dist       int
	generation uint32
}

func (a Node) Equal(b Node) bool {
	return a.Planet.Id == b.Planet.Id
}

type Graph []Node

type edge struct {
	from uint32
	to   uint32
	dist float64
}

type Edges []edge

func (a Edges) Len() int      { return len(a) }
func (a Edges) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a Edges) Less(i, j int) bool {
	return a[i].dist < a[j].dist
}

var generation uint32 = 0

func NewGraph(planets map[uint32]*pb.Planet) Graph {
	elems := make([]Node, len(planets))
	for i, planet := range planets {
		elems[i].Planet = planet
	}

	return elems
}

func (g Graph) Visit(root Node, f func(n Node) bool) {
	generation += 2
	queue := make([]*Node, 0)

	queue = append(queue, &g[root.Planet.Id])
	for len(queue) > 0 {
		elem := queue[0]
		queue = queue[1:]

		cont := f(g[elem.Planet.Id])
		elem.generation = generation
		if !cont {
			return
		}

		// visit children
		for _, child := range elem.Planet.Connected {
			if g[child].generation == generation {
				// already visited
				continue
			}
			queue = append(queue, &g[child])

			if g[child].generation < generation-1 || // never seen this child
				g[child].dist > elem.dist+1 { // the current path is shorter than the old one

				// set distance and paht
				g[child].dist = elem.dist + 1
				g[child].former = elem.Planet.Id
				// mark as seen
				g[child].generation = generation - 1
			}
		}
	}
}

func (g Graph) ShortestPath(root uint32, target uint32, graph_is_prepared bool) []uint32 {
	if ! graph_is_prepared {
		g.Visit(g[root], func(n Node) bool {
			return !g[target].Equal(n)
		})
	}

	path := make([]uint32, 0)

	cur := g[target]
	path = append(path, cur.Planet.Id)
	for !cur.Equal(g[root]) {
		path = append(path, cur.former)
		cur = g[cur.former]
		if len(path) > 500 {
			fmt.Println(root, g[root].Planet.Connected)
			fmt.Println(target, g[target].Planet.Connected)
			fmt.Println(path)
			for i := 0; i < 5; i++ {
				fmt.Println("  ", path[i], g[path[i]].Planet.Connected, "dist", g[path[i]].dist, "former", g[path[i]].former, "gen", g[path[i]].generation)
			}
			log.Panicf("Unintended behavior when checking path from %d to %d ", root, target, path)
		}
	}

	//revert path
	for i := 0; i < len(path)/2; i++ {
		path[i], path[len(path)-1-i ] = path[len(path)-1-i ], path[i]
	}

	return path
}

func (g Graph) GraphSize(root Node) int {

	count := 0

	g.Visit(g[root.Planet.Id], func(planet Node) bool {
		count++
		return true
	})

	return count
}

func (g Graph) HasCycle(root *pb.Planet) bool {
	generation++

	// remember which edges we already traveled
	usedEdges := make(map[edge]bool)

	// queue for bfs
	queue := make([]*Node, 0)
	queue = append(queue, &g[root.Id])

	for len(queue) > 0 {

		elem := queue[0]
		queue = queue[1:]

		// already visited
		if elem.generation == generation {
			return true
		}

		elem.generation = generation

		// visit children
		for _, child := range elem.Planet.Connected {
			e := edge{elem.Planet.Id, child, 0}
			// we come by this edge to this node therefore we can ignore it
			if usedEdges[e] {
				continue
			}
			usedEdges[edge{child, elem.Planet.Id, 0}] = true
			queue = append(queue, &g[child])
		}
	}

	return false
}
