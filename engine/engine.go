package engine

import (
	pb "github.com/eqinox76/RiseAndFallOfEmpires/proto"
	"github.com/eqinox76/RiseAndFallOfEmpires/state"
	"math/rand"
)

func Step(space *pb.Space) {
	for _, planet := range space.Planets {
		if rand.Float32() <= planet.Control {
			state.CreateShip(space, planet)
		}
	}
	return
}
