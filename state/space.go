package state

import (
	pb "github.com/eqinox76/RiseAndFallOfEmpires/proto"
	v "github.com/eqinox76/RiseAndFallOfEmpires/vector"
	"github.com/golang/protobuf/proto"
	"math/rand"
	"time"
	"github.com/dhconnelly/rtreego"

	"encoding/binary"
)

type Space struct{
	pb.Space
	PlanetTree *rtreego.Rtree
}

type PlanetPos struct {
	*pb.Planet
}

func (p PlanetPos) Bounds() *rtreego.Rect{
	return rtreego.Point{float64(p.PosX), float64(p.PosY)}.ToRect(0)
}

func NewSpace() Space {
	rand.Seed(time.Now().UTC().UnixNano())
	space := Space{
		Space: pb.Space{
			Width: 1000,
			Height: 500,
		},
		PlanetTree: rtreego.NewTree(2, 32, 64),
	}

	// add planets
	for i := uint32(0); i < 25; i++ {
		CreateNewPlanet(&space)
	}

	// add lanes between planets
	for _, planet := range space.Planets{
		nn := space.PlanetTree.NearestNeighbors(2, rtreego.Point{float64(planet.PosX), float64(planet.PosY)})
		ppos, ok := nn[1].(*PlanetPos)
		if !ok{
			panic(nn[1])
		}
		planet.Connected = append(planet.Connected, ppos.Id)
	}

	return space
}

func CreateShip(space *Space, planet *pb.Planet) *pb.Ship {
	var id uint64 = 0
	if space.Ships != nil {
		id = space.Ships[len(space.Ships)-1].Id
		id++
	}

	s := pb.Ship{
		Id: id,
	}

	space.Ships = append(space.Ships, &s)
	planet.Orbiting = append(planet.Orbiting, s.Id)
	return &s
}

func asVec(planet *pb.Planet) v.Vec {
	return v.Vec{float64(planet.PosX), float64(planet.PosY)}
}

func CreateNewPlanet(space *Space) *pb.Planet {
	var id uint32 = 0
	if space.Planets != nil {
		id = space.Planets[len(space.Planets)-1].Id
		id++
	}

	var x, y uint32
	for true {
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

			if vec.Dist(asVec(planet)) < 40 {
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
