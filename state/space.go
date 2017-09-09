package state

import (
	pb "github.com/eqinox76/RiseAndFallOfEmpires/proto"
	v "github.com/eqinox76/RiseAndFallOfEmpires/vector"
	"github.com/eqinox76/RiseAndFallOfEmpires/util"
	"github.com/golang/protobuf/proto"
	"math/rand"
	"github.com/dhconnelly/rtreego"

	"encoding/binary"
	"sort"
	"fmt"
	"math"
)

type Space struct {
	pb.Space
	PlanetTree *rtreego.Rtree
	Graph      Graph
	freeColors []string
}

func (space *Space) RemoveShip(ship *pb.Ship) {
	// remove from global ships
	space.Ships[ship.Id] = nil

	// remove from planet
	switch x := ship.GetPosition().(type) {
	case *pb.Ship_Orbiting:
		planet := space.Planets[x.Orbiting]
		util.RemoveUint64(&planet.Orbiting, ship.Id)
	default:
		panic(fmt.Sprintf("A destroyed ship is not orbiting a planet! %T", x))
	}

	// remove from empire
	empire := space.Empires[ship.Empire]
	util.RemoveUint64(&empire.Ships, ship.Id)
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
			Width:  1400,
			Height: 800,
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
		"darkgrey",
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
		"darkslategrey",
		"darkturquoise",
		"darkviolet",
		"deeppink",
		"deepskyblue",
		"dimgray",
		"dimgrey",
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
		"grey",
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
		"lightgrey",
		"lightpink",
		"lightsalmon",
		"lightseagreen",
		"lightskyblue",
		"lightslategray(Hex3)",
		"lightslategrey(Hex3)",
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
		"slategrey",
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
	space.freeColors = append(space.freeColors,neutralEmpire.Color)
	neutralEmpire.Color = "grey"
	neutralEmpire.Passive = true

	// add planets
	for i := uint32(0); i < 100; i++ {
		space.CreateNewPlanet(neutralEmpire)
	}

	// add empire start planets
	for ; empires > 0; empires-- {
		e := space.CreateEmpire()
		p := space.Planets[rand.Intn(len(space.Planets))]
		p.Empire = e.Id
		e.Planets = append(e.Planets, p.Id)
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
	// thereafter some planets are randomly connected with their n-nrearest neighbors but with decreasing likelyhood
	root := space.Planets[0]
	for _, edge := range edges {
		if space.Graph.GraphSize(root) == len(space.Planets) {
			// done
			break;
		}

		// add edge
		from := space.Planets[edge.from]
		to := space.Planets[edge.to]
		from.Connected = append(from.Connected, to.Id)
		to.Connected = append(to.Connected, from.Id)

		// check if cycling
		if space.Graph.HasCycle(from) {
			from.Connected = from.Connected[: len(from.Connected)-1]
			to.Connected = to.Connected[: len(to.Connected)-1]
		}
	}

	max_connections := int(math.Pow(float64(len(space.Planets)), 1.))
	for size := 1; size < max_connections; size++ {
		for _, planet := range space.Planets {
			if len(planet.Connected) == size {
				if rand.Intn(4) <= size {
					continue
				}

				nn := space.PlanetTree.NearestNeighbors(size+2, rtreego.Point{float64(planet.PosX), float64(planet.PosY)})

				to, ok := nn[size+1].(*PlanetPos)
				if ! ok {
					panic(nn[size+1])
				}

				planet.Connected = append(planet.Connected, to.Id)
				to.Connected = append(to.Connected, planet.Id)
			}
		}
	}

	// add neutral fleets
	for _, planet := range space.Planets{
		if planet.Empire == 0{
			ships := rand.Intn(10)
			for ; ships > 0; ships--{
				space.CreateShip(planet, space.Empires[0])
			}
		}
	}

	return space
}

func (space *Space) CreateShip(planet *pb.Planet, empire *pb.Empire) *pb.Ship {
	// TODO this creates a id from the last ship. but we may have already destroyed ships and therefore other open ids
	var id uint64 = 0
	if space.Ships != nil {
		id = space.Ships[len(space.Ships)-1].Id
		id++
	}

	s := pb.Ship{
		Id:     id,
		Empire: empire.Id,
		Position: &pb.Ship_Orbiting{
			Orbiting: planet.Id,
		},
	}

	space.Ships = append(space.Ships, &s)
	planet.Orbiting = append(planet.Orbiting, s.Id)
	empire.Ships = append(empire.Ships, s.Id)
	return &s
}

func (space *Space) CreateEmpire() *pb.Empire {
	var id uint32 = 0
	if space.Empires != nil {
		id = space.Empires[len(space.Empires)-1].Id
		id++
	}

	c := rand.Intn(len(space.freeColors))
	color := space.freeColors[c]
	space.freeColors[c] = space.freeColors[len(space.freeColors)-1]
	space.freeColors = space.freeColors[: len(space.freeColors)-1]

	e := pb.Empire{
		Id:    id,
		Color: color,
	}

	space.Empires = append(space.Empires, &e)
	return &e
}

func asVec(planet *pb.Planet) v.Vec {
	return v.Vec{float64(planet.PosX), float64(planet.PosY)}
}

func (space *Space) CreateNewPlanet(empire *pb.Empire) *pb.Planet {
	var id uint32 = 0
	if space.Planets != nil {
		id = space.Planets[len(space.Planets)-1].Id
		id++
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
		Id:      id,
		PosX:    x,
		PosY:    y,
		Control: rand.Float32(),
		Empire:  empire.Id,
	}

	space.Planets = append(space.Planets, &planet)

	space.PlanetTree.Insert(&PlanetPos{&planet})
	return &planet
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
