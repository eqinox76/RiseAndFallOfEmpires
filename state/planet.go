package state

import (
	pb "github.com/eqinox76/RiseAndFallOfEmpires/proto"
)

type Planet struct {
	id       uint32
	orbiting []*Ship
	Control  float32
	Empire   uint32
}

func (planet *Planet) Serialize() *pb.Planet {
	p := pb.Planet{
		Id:      planet.id,
		Control: planet.Control,
		Empire:  planet.Empire,
	}

	for _, s := range planet.orbiting {
		p.Orbiting = append(p.Orbiting, s.Id)
	}

	return &p
}