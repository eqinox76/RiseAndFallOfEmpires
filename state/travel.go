package state

import (
	pb "github.com/eqinox76/RiseAndFallOfEmpires/proto"
)

// only needed to add a marker to easily stop bfs
type node struct{
	planet *pb.Planet
	former *pb.Planet
	generation uint32
}

type Graph []node

type edge struct{
	from uint32
	to uint32
	dist float64
}

type Edges []edge

func (a Edges) Len() int           { return len(a) }
func (a Edges) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Edges) Less(i, j int) bool {
	return a[i].dist < a[j].dist
}

var  generation uint32 = 0

func NewGraph(planets map[uint32]*pb.Planet) Graph {
	elems := make([]node, len(planets))
	for i, planet := range planets {
		elems[i].planet = planet
	}

	return elems
}

func (elems Graph) Visit(root node, f func(planet *pb.Planet) bool){
	generation++
	queue := make([]node, 0)

	queue = append(queue, elems[root.planet.Id])
	for _, elem := range queue{
		// already visited
		if elem.generation == generation{
			continue
		}

		cont := f(elem.planet)
		if !cont{
			return
		}
		elem.generation = generation

		// visit children
		for _, child := range elem.planet.Connected{
			queue = append(queue, elems[child])
		}
	}
}

func (elems Graph) GraphSize(root *pb.Planet) int{

	count := 0

	elems.Visit(elems[root.Id], func(planet *pb.Planet) bool{
		count++
		return true
	})

	return count
}

func (elems Graph) HasCycle(root *pb.Planet) bool{
	generation++

	// remember which edges we already traveled
	usedEdges := make(map[edge] bool)

	// queue for bfs
	queue := make([]*node, 0)
	queue = append(queue, &elems[root.Id])

	for len(queue) > 0 {

		elem := queue[0]
		queue = queue[1:]

		// already visited
		if elem.generation == generation{
			return true
		}

		elem.generation = generation

		// visit children
		for _, child := range elem.planet.Connected{
			e := edge{elem.planet.Id, child, 0}
			// we come by this edge to this node therefore we can ignore it
			if usedEdges[e]{
				continue
			}
			usedEdges[edge{child, elem.planet.Id, 0}] = true
			queue = append(queue, &elems[child])
		}
	}

	return false
}
