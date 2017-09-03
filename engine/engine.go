package engine

import (
	"github.com/eqinox76/RiseAndFallOfEmpires/state"
	"math/rand"
)

func Step(space *state.Space) {
	for _, planet := range space.Planets {
		if rand.Float32() <= planet.Control {
			state.CreateShip(space, planet)
		}
		if planet.Control != 1 {
			if planet.Control > 0.999{
				planet.Control = 1
			} else if planet.Control > 0.5 {
				planet.Control = planet.Control + (rand.Float32() * 0.1 * (1 - planet.Control))
			} else {
				planet.Control = planet.Control + (rand.Float32() * 0.1 * (planet.Control))
			}
		}
	}
	return
}
