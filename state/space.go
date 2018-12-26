package state

import (
	"math/rand"

	"github.com/dhconnelly/rtreego"
	v "github.com/eqinox76/RiseAndFallOfEmpires/vector"
	"math"
	"sort"
)

type Empire struct {
	Passive bool
	Planets []*Planet
	Fleets  []*Fleet
}

type Space struct {
	Fleets  []*Fleet
	Empires []*Empire
	Planets []*Planet

	Round int

	PlanetTree *rtreego.Rtree
	Graph      Graph
	Width      uint32
	Height     uint32
}

func (planet Planet) Bounds() *rtreego.Rect {
	return rtreego.Point{float64(planet.PosX), float64(planet.PosY)}.ToRect(0)
}

func asVec(planet *Planet) v.Vec {
	return v.Vec{
		X: float64(planet.PosX),
		Y: float64(planet.PosY),
	}
}

func EmptySpace() Space {
	space := Space{
		Width:      1600,
		Height:     900,
		PlanetTree: rtreego.NewTree(2, 32, 64),
	}

	return space
}

func NewSpace(empires int) Space {
	space := EmptySpace()

	neutralEmpire := space.CreateEmpire()
	neutralEmpire.Passive = true

	// add planets
	for i := uint32(0); i < 130; i++ {
		space.CreatePlanet(neutralEmpire)
	}

	// add empire start planets
	for ; empires > 0; empires-- {
		e := space.CreateEmpire()

		p := space.Planets[uint32(rand.Int()%(len(space.Planets)))]
		for p.Empire != neutralEmpire {
			p = space.Planets[uint32(rand.Int()%(len(space.Planets)))]
		}

		e.Planets = append(e.Planets, p)
		p.Empire = e
		p.Control = 1
		p.Production = 0
	}

	space.Graph = NewGraph(space.Planets)
	// add the shortest paths until we have all nodes connected and do not add edges which add a cycle
	var edges Edges
	for i, planet := range space.Planets {
		for toI, to := range space.Planets {
			if i >= toI {
				continue
			}
			edges = append(edges, edge{i, toI, asVec(planet).Dist(asVec(to))})
		}
	}

	// sort the edged
	sort.Sort(edges)

	// add lanes between planets
	// this is done by adding all edged for a minimal spanning tree based on distance
	root := space.Graph.nodes[0]
	for _, edge := range edges {
		if space.Graph.GraphSize(root) == len(space.Planets) {
			// done
			break
		}

		// add edge
		from := space.Planets[edge.from]
		to := space.Planets[edge.to]
		from.Connected = append(from.Connected, to)
		to.Connected = append(to.Connected, from)

		// check if cycling
		if space.Graph.HasCycle(from) {
			from.Connected = from.Connected[:len(from.Connected)-1]
			to.Connected = to.Connected[:len(to.Connected)-1]
		}
	}

	// thereafter some planets are randomly connected with their n-nrearest neighbors but with decreasing likelyhood
	max_connections := int(math.Pow(float64(len(space.Planets)), 1./4.))
	for size := 1; size < max_connections; size++ {
		for _, planet := range space.Planets {
			if rand.Intn(max_connections) < size {
				continue
			}

			nn := space.PlanetTree.NearestNeighbors(size+2, rtreego.Point{float64(planet.PosX), float64(planet.PosY)})

			to, ok := nn[size+1].(*Planet)
			if !ok {
				panic(nn[size+1])
			}

			// check if this edge already exists
			exists := false
			for _, id := range planet.Connected {
				if id == to {
					exists = true
					break
				}
			}

			if !exists {
				planet.Connected = append(planet.Connected, to)
				to.Connected = append(to.Connected, planet)
			}
		}
	}

	// add neutral fleets
	for _, planet := range space.Planets {
		if planet.Empire.Passive {
			fleet := space.CreateFleet(planet, planet.Empire)
			fleet.LightSquads = 2 + rand.Intn(8)
		}
	}

	return space
}

func (space *Space) CreateEmpire() *Empire {

	e := Empire{}

	space.Empires = append(space.Empires, &e)
	return &e
}

func (space *Space) CreateFleet(planet *Planet, empire *Empire) *Fleet {
	f := Fleet{
		Empire:   empire,
		Position: planet,
	}
	planet.Fleets = append(planet.Fleets, &f)
	empire.Fleets = append(empire.Fleets, &f)
	space.Fleets = append(space.Fleets, &f)
	return &f
}

func (space *Space) CreatePlanet(empire *Empire) *Planet {
	var x, y uint32
	for {
		x, y = rand.Uint32()%space.Width, rand.Uint32()%space.Height
		valid := true
		if x < 50 || x > space.Width-50 {
			continue
		}

		if y < 50 || y > space.Height-50 {
			continue
		}

		vec := v.Vec{
			X: float64(x),
			Y: float64(y),
		}
		for _, planet := range space.Planets {

			if vec.Dist(asVec(planet)) < 50 {
				valid = false
				break
			}
		}

		if valid {
			break
		}
	}

	planet := Planet{
		PosX:    x,
		PosY:    y,
		Control: rand.Float32(),
		Empire:  empire,
	}

	empire.Planets = append(empire.Planets, &planet)
	space.Planets = append(space.Planets, &planet)
	space.PlanetTree.Insert(&planet)

	return &planet
}

func (space *Space) DestroyFleet(fleet *Fleet){
	rm := func (fleets []*Fleet,f *Fleet) []*Fleet{
		for pos, val := range fleets{
			if val == f{
				fleets[pos] = fleets[len(fleets) - 1]
				return fleets[:len(fleets) - 1]
			}
		}
		return fleets
	}

	space.Fleets =rm(space.Fleets, fleet)
	fleet.Empire.Fleets = rm(fleet.Empire.Fleets, fleet)
	fleet.Position.Fleets = rm(fleet.Position.Fleets, fleet)
}

func (space *Space) Won() bool {
	return len(space.Empires) == 2
}
