package model

import (
	pb "github.com/eqinox76/RiseAndFallOfEmpires/proto"
	"github.com/golang/protobuf/proto"
	"math/rand"
)

type Space struct {
	internal     pb.Space
	nextShipId   uint64
	nextPlanetId uint32
}

type Planet struct {
	internal pb.Planet
}
type Ship struct{}

func NewSpace() Space {
	space := Space{
		internal: pb.Space{},
	}

	space.nextPlanetId = 0
	space.nextShipId = 0

	for i := uint32(0); i < 100; i++ {
		space.createPlanet(i);
	}

	return space
}

func (space *Space) createShip(planet *Planet) Ship {
	return Ship{}
}

func (space *Space) createPlanet(empire uint32) Planet {
	id := space.nextPlanetId
	space.nextPlanetId++
	planet := Planet{
		pb.Planet{
			Id:      id,
			Control: rand.Float32(),
			Empire:  empire,
		},
	}

	space.internal.Planets = append(space.internal.Planets, &planet.internal)

	return planet
}

func (space *Space) Serialize() ([]byte, error) {
	return proto.Marshal(&space.internal)
}
