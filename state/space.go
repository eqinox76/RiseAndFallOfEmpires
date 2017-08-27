package state

import (
	pb "github.com/eqinox76/RiseAndFallOfEmpires/proto"
	"github.com/golang/protobuf/proto"
	"math/rand"
	"github.com/pkg/errors"
)

type Space struct {
	Planets      []*Planet
	Ships        []*Ship
	nextShipId   uint64
	nextPlanetId uint32
}

func NewSpace() Space {
	space := Space{
		nextShipId:   0,
		nextPlanetId: 0,
	}

	for i := uint32(0); i < 10; i++ {
		space.CreatePlanet(i);
	}

	return space
}

func (space *Space) CreateShip(planet *Planet) Ship {
	id := space.nextShipId
	space.nextShipId++

	s := Ship{
		Id: id,
	}
	space.Ships = append(space.Ships, &s)
	planet.orbiting = append(planet.orbiting, &s)
	return s
}

func (space *Space) CreatePlanet(empire uint32) *Planet {
	id := space.nextPlanetId
	space.nextPlanetId++
	planet := Planet{
		id:      id,
		Control: rand.Float32(),
		Empire:  empire,
	}

	space.Planets = append(space.Planets, &planet)

	return &planet
}

func (space *Space) Serialize() ([]byte, error) {
	out := pb.Space{}
	for _, planet := range space.Planets {
		out.Planets = append(out.Planets, planet.Serialize())
	}
	return proto.Marshal(&out)
}

func Deserialize(data *[]byte) (*Space, error) {
	in := pb.Space{}
	err := proto.Unmarshal(*data, &in)
	if err != nil {
		return nil, errors.Wrap(err, "Could not deserialize.")
	}
	space := NewSpace()

	for _, planet := range in.Planets {
		p := space.CreatePlanet(planet.Empire)
		for _, _ = range planet.Orbiting{
			space.CreateShip(p)
		}
	}

	return &space, nil
}
