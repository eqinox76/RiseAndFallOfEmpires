package model

import (
	pb "github.com/eqinox76/RiseAndFallOfEmpires/proto"
	"github.com/golang/protobuf/proto"
	"math/rand"
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

func (space *Space) CreatePlanet(empire uint32) Planet {
	id := space.nextPlanetId
	space.nextPlanetId++
	planet := Planet{
		id: id,
		Control: rand.Float32(),
	}

	space.Planets = append(space.Planets, &planet)

	return planet
}

func (space *Space) Serialize() ([]byte, error) {
	out := pb.Space{}
	for _, planet := range space.Planets {
		out.Planets = append(out.Planets, planet.Serialize())
	}
	return proto.Marshal(&out)
}

