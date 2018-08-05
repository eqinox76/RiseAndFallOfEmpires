package state

import (
	"math/rand"

	"github.com/dhconnelly/rtreego"
	pb "github.com/eqinox76/RiseAndFallOfEmpires/proto"
	"github.com/eqinox76/RiseAndFallOfEmpires/util"
	v "github.com/eqinox76/RiseAndFallOfEmpires/vector"
	"github.com/golang/protobuf/proto"

	"encoding/binary"
	"fmt"
	"math"
	"sort"
)

type Space struct {
	pb.Space
	PlanetTree *rtreego.Rtree
	Graph      Graph
	freeColors []string
}

// MoveFleet for now lets just teleport ships
func (space *Space) MoveFleet(fleetid uint32, start uint32, destination uint32) {
	space.Fleets[fleetid].Position = destination

	space.Planets[start].Fleets = util.Removeu32(space.Planets[start].Fleets, fleetid)
	space.Planets[start].Fleets = append(space.Planets[destination].Fleets, fleetid)
}

// RemoveShip removes a ship. throws if none is available
func (space *Space) RemoveShip(shiptype uint32, fleetid uint32) {
	available := space.Fleets[fleetid].Ships[shiptype]
	if available == 0 {
		panic(fmt.Sprintf("Cannot remove because fleet %T has no %T!", fleetid, shiptype))
	} else {
		space.Fleets[fleetid].Ships[shiptype]--
	}
}

type PlanetPos struct {
	*pb.Planet
}

func (p PlanetPos) Bounds() *rtreego.Rect {
	return rtreego.Point{float64(p.PosX), float64(p.PosY)}.ToRect(0)
}

func EmptySpace() Space {
	space := Space{
		Space: pb.Space{
			Width:   1000,
			Height:  600,
			Fleets:  nil,
			Planets: make(map[uint32]*pb.Planet),
			Empires: make(map[uint32]*pb.Empire),
		},
		PlanetTree: rtreego.NewTree(2, 32, 64),
	}

	space.freeColors = []string{
		"aliceblue",
		"antiquewhite",
		"aqua",
		"aquamarine",
		"azure",
		"beige",
		"bisque",
		"black",
		"blanchedalmond",
		"blue",
		"blueviolet",
		"brown",
		"burlywood",
		"cadetblue",
		"chartreuse",
		"chocolate",
		"coral",
		"cornflowerblue",
		"cornsilk",
		"crimson",
		"cyan",
		"darkblue",
		"darkcyan",
		"darkgoldenrod",
		"darkgray",
		"darkgreen",
		"darkkhaki",
		"darkmagenta",
		"darkolivegreen",
		"darkorange",
		"darkorchid",
		"darkred",
		"darksalmon",
		"darkseagreen",
		"darkslateblue",
		"darkslategray",
		"darkturquoise",
		"darkviolet",
		"deeppink",
		"deepskyblue",
		"dimgray",
		"dodgerblue",
		"firebrick",
		"floralwhite",
		"forestgreen",
		"fuchsia",
		"gainsboro",
		"ghostwhite",
		"gold",
		"goldenrod",
		"green",
		"greenyellow",
		"honeydew",
		"hotpink",
		"indianred",
		"indigo",
		"ivory",
		"khaki",
		"lavender",
		"lavenderblush",
		"lawngreen",
		"lemonchiffon",
		"lightblue",
		"lightcoral",
		"lightcyan",
		"lightgoldenrodyellow",
		"lightgray",
		"lightgreen",
		"lightpink",
		"lightsalmon",
		"lightseagreen",
		"lightskyblue",
		"lightslategray",
		"lightsteelblue",
		"lightyellow",
		"lime",
		"limegreen",
		"linen",
		"magenta",
		"maroon",
		"mediumaquamarine",
		"mediumblue",
		"mediumorchid",
		"mediumpurple",
		"mediumseagreen",
		"mediumslateblue",
		"mediumspringgreen",
		"mediumturquoise",
		"mediumvioletred",
		"midnightblue",
		"mintcream",
		"mistyrose",
		"moccasin",
		"navajowhite",
		"navy",
		"oldlace",
		"olive",
		"olivedrab",
		"orange",
		"orangered",
		"orchid",
		"palegoldenrod",
		"palegreen",
		"paleturquoise",
		"palevioletred",
		"papayawhip",
		"peachpuff",
		"peru",
		"pink",
		"plum",
		"powderblue",
		"purple",
		"red",
		"rosybrown",
		"royalblue",
		"saddlebrown",
		"salmon",
		"sandybrown",
		"seagreen",
		"seashell",
		"sienna",
		"silver",
		"skyblue",
		"slateblue",
		"slategray",
		"snow",
		"springgreen",
		"steelblue",
		"tan",
		"teal",
		"thistle",
		"tomato",
		"turquoise",
		"violet",
		"wheat",
		"white",
		"whitesmoke",
		"yellow",
		"yellowgreen}"}

	return space
}

func NewSpace(empires int) Space {
	space := EmptySpace()

	neutralEmpire := space.CreateEmpire()
	space.freeColors = append(space.freeColors, neutralEmpire.Color)
	neutralEmpire.Color = "grey"
	neutralEmpire.Passive = true

	// add planets
	for i := uint32(0); i < 100; i++ {
		space.CreatePlanet(neutralEmpire)
	}

	// add empire start planets
	for ; empires > 0; empires-- {
		e := space.CreateEmpire()

		p := space.Planets[uint32(rand.Int()%(len(space.Planets)))]
		for p.Empire != 0 {
			p = space.Planets[uint32(rand.Int()%(len(space.Planets)))]
		}

		p.Empire = e.Id
		e.Planets[p.Id] = true
		p.Control = 1
		p.Production = 0
	}

	space.Graph = NewGraph(space.Planets)
	// add the shortest paths until we have all nodes connected and do not add edges which add a cycle
	var edges Edges
	for _, planet := range space.Planets {
		for _, to := range space.Planets {
			if planet.Id >= to.Id {
				continue
			}
			edges = append(edges, edge{planet.Id, to.Id, asVec(planet).Dist(asVec(to))})
		}
	}

	// sort the edged
	sort.Sort(edges)

	// add lanes between planets
	// this is done by adding all edged for a minimal spanning tree based on distance
	root := space.Graph.nodes[space.Planets[0].Id]
	for _, edge := range edges {
		if space.Graph.GraphSize(root) == len(space.Planets) {
			// done
			break
		}

		// add edge
		from := space.Planets[edge.from]
		to := space.Planets[edge.to]
		from.Connected = append(from.Connected, to.Id)
		to.Connected = append(to.Connected, from.Id)

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

			to, ok := nn[size+1].(*PlanetPos)
			if !ok {
				panic(nn[size+1])
			}

			// check if this edge already exists
			exists := false
			for _, id := range planet.Connected {
				if id == to.Id {
					exists = true
					break
				}
			}

			if !exists {
				planet.Connected = append(planet.Connected, to.Id)
				to.Connected = append(to.Connected, planet.Id)
			}
		}
	}

	// add neutral fleets
	for _, planet := range space.Planets {
		if planet.Empire == 0 {
			ships := 2 + rand.Intn(8)
			for ; ships > 0; ships-- {
				space.CreateShip(planet, space.Empires[0])
			}
		}
	}

	return space
}

func (space *Space) CreateShip(planet *pb.Planet, empire *pb.Empire) *pb.Ship {
	id := uint64(len(space.Ships))
	_, contained := space.Ships[id]
	for contained {
		id++
		_, contained = space.Ships[id]
	}

	s := pb.Ship{
		Id:     id,
		Empire: empire.Id,
		Position: &pb.Ship_Orbiting{
			Orbiting: planet.Id,
		},
	}

	space.Ships[id] = &s
	planet.Orbiting[s.Id] = true
	empire.Ships[s.Id] = true
	return &s
}

func (space *Space) CreateEmpire() *pb.Empire {
	id := uint32(len(space.Empires))
	_, contained := space.Empires[id]
	for ; contained; id++ {
		_, contained = space.Empires[id]
	}

	c := rand.Intn(len(space.freeColors))
	color := space.freeColors[c]
	space.freeColors[c] = space.freeColors[len(space.freeColors)-1]
	space.freeColors = space.freeColors[:len(space.freeColors)-1]

	e := pb.Empire{
		Id:      id,
		Color:   color,
		Ships:   make(map[uint64]bool),
		Planets: make(map[uint32]bool),
	}

	space.Empires[e.Id] = &e
	return &e
}

func asVec(planet *pb.Planet) v.Vec {
	return v.Vec{float64(planet.PosX), float64(planet.PosY)}
}

func (space *Space) CreatePlanet(empire *pb.Empire) *pb.Planet {
	id := uint32(len(space.Planets))
	_, contained := space.Planets[id]
	for ; contained; id++ {
		_, contained = space.Planets[id]
	}

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

		vec := v.Vec{float64(x), float64(y)}
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

	planet := pb.Planet{
		Id:       id,
		PosX:     x,
		PosY:     y,
		Control:  rand.Float32(),
		Empire:   empire.Id,
		Orbiting: make(map[uint64]bool),
	}

	empire.Planets[planet.Id] = true
	space.Planets[planet.Id] = &planet
	space.PlanetTree.Insert(&PlanetPos{&planet})

	return &planet
}

func GetFleets(global_ships map[uint64]*pb.Ship, planet *pb.Planet) map[uint32][]*pb.Ship {
	fleets := make(map[uint32][]*pb.Ship)

	for id, _ := range planet.Orbiting {
		s := global_ships[id]
		_, ok := fleets[s.Empire]
		if !ok {
			fleets[s.Empire] = []*pb.Ship{}
		}

		fleets[s.Empire] = append(fleets[s.Empire], s)
	}

	return fleets
}

func (space *Space) Won() bool {
	return len(space.Empires) == 2
}

func Serialize(space *Space) ([]byte, error) {

	data, err := proto.Marshal(&space.Space)
	if err != nil {
		return nil, err
	}

	length := make([]byte, 4)
	binary.LittleEndian.PutUint32(length, uint32(len(data)))

	return append(length, data...), nil
}
