package engine

import (
	"github.com/eqinox76/RiseAndFallOfEmpires/state"
	"math/rand"
	"math"
)

func Step(space *state.Space) {
	for _, planet := range space.Planets {
		if rand.Float32() <= planet.Control {
			state.CreateShip(space, planet)
		}
		if planet.Control != 1 {
			if planet.Control > 0.999 {
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

func Fight(shipsA uint32, shipsB uint32) (uint32, uint32) {
	a := shipsA - destroyed(shipsB, shipsA)
	b := shipsB - destroyed(shipsA, shipsB)
	return a, b
}

func destroyed(ships uint32, target uint32) uint32 {
	const prob float64 = 0.1
	var deviation float64 = math.Sqrt(float64(ships) * prob * (1 - prob))
	var destroyed float64 = rand.NormFloat64()*deviation + (float64(ships) * prob)

	if destroyed < 0 {
		return 0
	} else if target < uint32(destroyed) {
		return target
	} else {
		return uint32(destroyed)
	}
}
