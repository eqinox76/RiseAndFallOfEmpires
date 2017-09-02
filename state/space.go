package state

import (
	pb "github.com/eqinox76/RiseAndFallOfEmpires/proto"
	"github.com/golang/protobuf/proto"
	"math/rand"
	"time"

	"encoding/binary"
)

func NewSpace() pb.Space {
	rand.Seed(time.Now().UTC().UnixNano())
	space := pb.Space{
		Width: 1400,
		Height: 500,
	}

	for i := uint32(0); i < 150; i++ {
		CreateNewPlanet(&space);
	}

	return space
}

func CreateShip(space *pb.Space, planet *pb.Planet) *pb.Ship {
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

func CreateNewPlanet(space *pb.Space) *pb.Planet {
	var id uint32 = 0
	if space.Planets != nil {
		id = space.Planets[len(space.Planets)-1].Id
		id++
	}

	var x,y uint32
	for true{
		x, y = rand.Uint32() % space.Width, rand.Uint32() % space.Height
		if x < 50 || x > space.Width - 50{
			continue
		}

		if y < 50 || y > space.Height - 50{
			continue
		}

		for _, planet := range space.Planets{
			if x == planet.PosX || y == planet.PosY {
				continue
			}
		}
		break
	}
	planet := pb.Planet{
		Id: id,
		PosX:  x,
		PosY:  y,
		Control: rand.Float32(),
	}

	space.Planets = append(space.Planets, &planet)

	return &planet
}

func Serialize(space *pb.Space) ([]byte, error) {

	data, err := proto.Marshal(space)
	if err != nil {
		return nil, err
	}

	length := make([]byte, 4)
	binary.LittleEndian.PutUint32(length, uint32(len(data)))

	return append(length, data...), nil
}
