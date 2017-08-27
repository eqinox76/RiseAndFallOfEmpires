package state

import pb "github.com/eqinox76/RiseAndFallOfEmpires/proto"

type Ship struct {
	Id uint64
}

func (ship *Ship) Serialize() pb.Ship {
	s := pb.Ship{
		Id: ship.Id,
	}
	return s
}
