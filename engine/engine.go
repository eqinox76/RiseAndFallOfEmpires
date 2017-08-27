package engine

import (
	"github.com/eqinox76/RiseAndFallOfEmpires/state"
	"math/rand"
)

func Step(space *state.Space) {
	for _, planet := range space.Planets {
		if rand.Float32() <= planet.Control {
			space.CreateShip(planet)
		}
	}
	return
}
