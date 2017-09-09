package state

import (
	pb "github.com/eqinox76/RiseAndFallOfEmpires/proto"
	v "github.com/eqinox76/RiseAndFallOfEmpires/vector"
	"github.com/eqinox76/RiseAndFallOfEmpires/util"
	"github.com/golang/protobuf/proto"
	"math/rand"
	"time"
	"github.com/dhconnelly/rtreego"

	"encoding/binary"
	"sort"
	"fmt"
)

type Space struct {
	pb.Space
	PlanetTree *rtreego.Rtree
	Graph      Graph
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
	rand.Seed(time.Now().UTC().UnixNano())
	space := Space{
		Space: pb.Space{
			Width:  1000,
			Height: 500,
		},
		PlanetTree: rtreego.NewTree(2, 32, 64),
	}

	return space
}

func NewSpace() Space {
	space := EmptySpace()

	// add planets
	for i := uint32(0); i < 25; i++ {
		space.CreateNewPlanet()
	}

	space.Graph = NewGraph(space.Planets)
	// add the shortes paths until we have all nodes connected and do not add edges which add a cycle
	var edges Edges
	for _, planet := range space.Planets{
		for _, to := range space.Planets{
			if planet.Id >= to.Id{
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
	for _, edge := range edges{
		if space.Graph.GraphSize(root) == len(space.Planets){
			// done
			break;
		}

		// add edge
		from := space.Planets[edge.from]
		to := space.Planets[edge.to]
		from.Connected = append(from.Connected, to.Id)
		to.Connected = append(to.Connected, from.Id)

		// check if cycling
		if space.Graph.HasCycle(from){
			from.Connected = from.Connected[: len(from.Connected) - 1]
			to.Connected = to.Connected[: len(to.Connected) - 1]
		}
	}

	for size := 1; size < 3; size++ {
		for _, planet := range space.Planets {
			if len(planet.Connected) == size {
				if rand.Intn(4) <= size{
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
		Id: id,
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

func (space *Space) CreateEmpire(color string) *pb.Empire{
	var id uint32 = 0
	if space.Empires != nil {
		id = space.Empires[len(space.Empires)-1].Id
		id++
	}

	e := pb.Empire{
		Id: id,
		Color: color,
	}

	space.Empires = append(space.Empires, &e)
	return &e
}

func asVec(planet *pb.Planet) v.Vec {
	return v.Vec{float64(planet.PosX), float64(planet.PosY)}
}

func (space *Space) CreateNewPlanet() *pb.Planet {
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
